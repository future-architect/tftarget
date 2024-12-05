/*
Copyright © 2023 Takafumi Miyanaga miya.org.0309@gmai.com
*/
package cmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
)

func Test_extractResource(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []string
	}{
		{
			name: "plan結果からリソース名を取得1",
			input: []byte(`
      Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the
      following symbols:
        + create

      Terraform will perform the following actions:

        # aws_iam_role_policy_attachment.local_poilcy_attachment will be created
        + resource "aws_iam_role_policy_attachment" "local_poilcy_attachment" {
            + id         = (known after apply)
            + policy_arn = "arn:aws:iam::000000000000:policy/local-iam-policy"
            + role       = "local-role"
          }

        # aws_lambda_function.local_lambda will be created
        + resource "aws_lambda_function" "local_lambda" {
            + architectures                  = (known after apply)
            + arn                            = (known after apply)
            + filename                       = "lambda.zip"
            + function_name                  = "local-lambda"
            + handler                        = "lambda"
            + id                             = (known after apply)
            + invoke_arn                     = (known after apply)
            + last_modified                  = (known after apply)
            + memory_size                    = 128
            + package_type                   = "Zip"
            + publish                        = false
            + qualified_arn                  = (known after apply)
            + qualified_invoke_arn           = (known after apply)
            + reserved_concurrent_executions = -1
            + role                           = "arn:aws:iam::000000000000:role/local-role"
            + runtime                        = "go1.x"
            + signing_job_arn                = (known after apply)
            + signing_profile_version_arn    = (known after apply)
            + skip_destroy                   = false
            + source_code_hash               = (known after apply)
            + source_code_size               = (known after apply)
            + tags_all                       = (known after apply)
            + timeout                        = 3
            + version                        = (known after apply)

            + ephemeral_storage {
                + size = (known after apply)
              }

            + tracing_config {
                + mode = (known after apply)
              }
          }
          # aws_iam_policy.local_policy will be destroyed
          - resource "aws_iam_policy" "local_policy" {
              - arn       = "arn:aws:iam::000000000000:policy/local-iam-policy" -> null
              - id        = "arn:aws:iam::000000000000:policy/local-iam-policy" -> null
              - name      = "local-iam-policy" -> null
              - path      = "/" -> null
              - policy    = jsonencode(
                    {
                      - Statement = [
                          - {
                              - Action   = [
                                  - "s3:PutObject",
                                  - "s3:PutBucketNotification",
                                  - "s3:ListBucket",
                                  - "s3:GetObject",
                                  - "s3:DeleteObject",
                                ]
                              - Effect   = "Allow"
                              - Resource = "*"
                              - Sid      = ""
                            },
                          - {
                              - Action   = [
                                  - "kinesis:GetShardIterator",
                                  - "kinesis:GetRecords",
                                  - "kinesis:DescribeStream",
                                ]
                              - Effect   = "Allow"
                              - Resource = "*"
                              - Sid      = ""
                            },
                        ]
                      - Version   = "2012-10-17"
                    }
                ) -> null
              - policy_id = "ALQL2Y1NHP2PZP7QD45PC" -> null
              - tags      = {} -> null
              - tags_all  = {} -> null
            }
          `),
			want: []string{"aws_iam_role_policy_attachment.local_poilcy_attachment will be created", "aws_lambda_function.local_lambda will be created", "aws_iam_policy.local_policy will be destroyed"},
		},
		{
			name: "# (xxx unchanged)をリソース名に含めない",
			input: []byte(`
      Terraform will perform the following actions:

      # aws_lambda_event_source_mapping.local_mapping will be updated in-place
      ~ resource "aws_lambda_event_source_mapping" "local_mapping" {
          ~ batch_size                         = 100 -> 130
            id                                 = "4a8c06dc-fec4-48be-8bce-8f8beb445a57"
            # (19 unchanged attributes hidden)
        }
          `),
			want: []string{"aws_lambda_event_source_mapping.local_mapping will be updated in-place"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractResource(tt.input, "")
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("extractResouce (-got +want):%s", diff)
			}
		})
	}
}

func Test_dropAction(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "リソース名のみを取得",
			input: []string{"aws_iam_role_policy_attachment.local_poilcy_attachment will be created", "aws_lambda_function.local_lambda will be created", "aws_iam_policy.local_policy will be destroyed", "aws_lambda_event_source_mapping.local_mapping will be updated in-place"},
			want:  []string{"aws_iam_role_policy_attachment.local_poilcy_attachment", "aws_lambda_function.local_lambda", "aws_iam_policy.local_policy", "aws_lambda_event_source_mapping.local_mapping"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dropAction(tt.input)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("dropAction (-got +want):%s", diff)
			}
		})
	}
}

func Test_slice2String(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{
			name:  "スライスをストリングに変換",
			input: []string{"aws_iam_role_policy_attachment.local_poilcy_attachment", "aws_lambda_function.local_lambda", "aws_iam_policy.local_policy", "aws_lambda_event_source_mapping.local_mapping", "aws_lambda_function.lambda_set[\"test\"]"},
			want:  "{'aws_iam_role_policy_attachment.local_poilcy_attachment','aws_lambda_function.local_lambda','aws_iam_policy.local_policy','aws_lambda_event_source_mapping.local_mapping','aws_lambda_function.lambda_set[\"test\"]'}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slice2String(tt.input)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("slice2String (-got +want):%s", diff)
			}
		})
	}
}

func Test_genTargetCmd(t *testing.T) {
	tests := []struct {
		name   string
		action string
		input  string
		want   string
		cmd    *cobra.Command
	}{
		{
			name:   "plan -target",
			cmd:    planCmd,
			action: "plan",
			input:  "{aws_iam_role_policy_attachment.local_poilcy_attachment,aws_lambda_function.local_lambda,aws_iam_policy.local_policy,aws_lambda_event_source_mapping.local_mapping}",
			want:   "terraform plan -target={aws_iam_role_policy_attachment.local_poilcy_attachment,aws_lambda_function.local_lambda,aws_iam_policy.local_policy,aws_lambda_event_source_mapping.local_mapping} --parallelism=10",
		},
		{
			name:   "apply -target",
			cmd:    applyCmd,
			action: "apply",
			input:  "{aws_iam_role_policy_attachment.local_poilcy_attachment,aws_lambda_function.local_lambda,aws_iam_policy.local_policy,aws_lambda_event_source_mapping.local_mapping}",
			want:   "terraform apply -target={aws_iam_role_policy_attachment.local_poilcy_attachment,aws_lambda_function.local_lambda,aws_iam_policy.local_policy,aws_lambda_event_source_mapping.local_mapping} --parallelism=10",
		},
		{
			name:   "destroy -target",
			cmd:    destroyCmd,
			action: "destroy",
			input:  "{aws_iam_role_policy_attachment.local_poilcy_attachment,aws_lambda_function.local_lambda,aws_iam_policy.local_policy,aws_lambda_event_source_mapping.local_mapping}",
			want:   "terraform destroy -target={aws_iam_role_policy_attachment.local_poilcy_attachment,aws_lambda_function.local_lambda,aws_iam_policy.local_policy,aws_lambda_event_source_mapping.local_mapping} --parallelism=10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The --executable flag is a persistent flag set by the root command, but because the default value of it
			// is not passed to the command because it was not executed through the root command, the flag is set here
			// as flag with the expected test value set as default value.
			tt.cmd.Flags().String("executable", "terraform", "")
			got := genTargetCmd(tt.cmd, tt.action, tt.input)
			if diff := cmp.Diff(got.String(), tt.want); diff != "" {
				t.Errorf("genTargetCmd (-got +want):%s", diff)
			}
		})
	}
}
