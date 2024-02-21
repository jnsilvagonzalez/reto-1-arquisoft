resource "aws_iam_role" "lambda_role" {
  name = "lambda_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_policy" "lambda_policy" {
  name        = "lambda_policy"
  path        = "/"
  description = "Lambda policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "sns:*", "xray:*", "dynamodb:*", "logs:*",
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}

resource "aws_iam_policy_attachment" "lambda_attach" {
  name       = "lambda-attachment"
  roles      = [aws_iam_role.lambda_role.name]
  policy_arn = aws_iam_policy.lambda_policy.arn
}

resource "aws_lambda_function" "emergency_lambda" {
  filename      = "${path.module}\\bootstrap.zip"
  function_name = "procesador_eventos_emergencia"
  role          = aws_iam_role.lambda_role.arn
  publish       = true
  handler       = "main"
  runtime       = "provided.al2023"
  timeout       = 2
}

# resource "aws_lambda_provisioned_concurrency_config" "lambda" {
#   function_name                     = aws_lambda_function.emergency_lambda.function_name
#   provisioned_concurrent_executions = 1
#   qualifier                         = aws_lambda_function.emergency_lambda.version
# }

# resource "aws_lambda_function" "emergency_lambda_python" {
#   filename      = "${path.module}\\lambda.zip"
#   function_name = "procesador_eventos_emergencia_python"
#   role          = aws_iam_role.lambda_role.arn
#   handler       = "lambda_function.lambda_handler"
#   runtime       = "python3.10"
#   timeout       = 2
# }

# resource "aws_lambda_permission" "sns_invoke_python" {
#   statement_id  = "AllowExecutionFromSNS"
#   action        = "lambda:InvokeFunction"
#   function_name = aws_lambda_function.emergency_lambda_python.function_name
#   principal     = "sns.amazonaws.com"
#   source_arn    = aws_sns_topic.emergency.arn
# }

# resource "aws_sns_topic_subscription" "lambda_python" {
#   topic_arn = aws_sns_topic.emergency.arn
#   protocol  = "lambda"
#   endpoint  = aws_lambda_function.emergency_lambda_python.arn
# }

resource "aws_lambda_permission" "sns_invoke" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.emergency_lambda.function_name
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.emergency.arn
}

resource "aws_sns_topic_subscription" "lambda" {
  topic_arn = aws_sns_topic.emergency.arn
  protocol  = "lambda"
  endpoint  = aws_lambda_function.emergency_lambda.arn
}


resource "aws_dynamodb_table" "basic-dynamodb-table" {
  name         = "Rules_table"
  billing_mode = "PAY_PER_REQUEST"

  hash_key  = "VehiculoID"
  range_key = "ReglaId"

  attribute {
    name = "VehiculoID"
    type = "S"
  }

  attribute {
    name = "ReglaId"
    type = "S"
  }
}
