resource "aws_security_group" "alb_http_access" {
  name   = "alb_http_access"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "ec2_http_access" {
  name   = "ec2_http_access"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port       = 8080
    to_port         = 8084
    protocol        = "tcp"
    security_groups = [aws_security_group.alb_http_access.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
resource "aws_lb" "receptor" {
  name               = "alb-receptor"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb_http_access.id]
  subnets            = [aws_subnet.az_1.id, aws_subnet.az_2.id]
}

resource "aws_lb_listener" "receptor" {
  load_balancer_arn = aws_lb.receptor.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.signals_receptor_1.arn
  }
}

resource "aws_lb_listener_rule" "emergency" {
  listener_arn = aws_lb_listener.receptor.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.emergency_receptor.arn
  }
  condition {
    path_pattern {
      values = ["/api/emergency"]
    }
  }
}

resource "aws_lb_listener_rule" "signals" {
  listener_arn = aws_lb_listener.receptor.arn
  priority     = 200

  action {
    type = "forward"
    forward {
      target_group {
        arn = aws_lb_target_group.signals_receptor_1.arn
      }
      target_group {
        arn = aws_lb_target_group.signals_receptor_2.arn
      }
      target_group {
        arn = aws_lb_target_group.signals_receptor_3.arn
      }
      target_group {
        arn = aws_lb_target_group.signals_receptor_4.arn
      }
    }
  }
  condition {
    path_pattern {
      values = ["/api/signal"]
    }
  }
}

resource "aws_lb_target_group" "signals_receptor_1" {
  name     = "signals-receptor-1"
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id
}

resource "aws_lb_target_group" "signals_receptor_2" {
  name     = "signals-receptor-2"
  port     = 8081
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id
}

resource "aws_lb_target_group" "signals_receptor_3" {
  name     = "signals-receptor-3"
  port     = 8082
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id
}

resource "aws_lb_target_group" "signals_receptor_4" {
  name     = "signals-receptor-4"
  port     = 8083
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id
}


resource "aws_lb_target_group" "emergency_receptor" {
  name     = "emergency-receptor"
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id
}

resource "aws_iam_role" "ec2_role" {
  name = "ec2_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_policy" "ec2_policy" {
  name        = "ec2_policy"
  path        = "/"
  description = "Ec2 policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "sns:*", "xray:*",
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}

resource "aws_iam_policy_attachment" "ec2_attach" {
  name       = "ec2-attachment"
  roles      = [aws_iam_role.ec2_role.name]
  policy_arn = aws_iam_policy.ec2_policy.arn
}

resource "aws_launch_template" "normal_receptor" {
  name     = "normal_receptor_launch_template"
  image_id = "ami-000807bf3875b8d6c"

  instance_market_options {
    market_type = "spot"
  }

  network_interfaces {
    security_groups             = [aws_security_group.ec2_http_access.id]
    associate_public_ip_address = true
  }

  instance_type = "m5.xlarge"
  user_data = base64encode(
    <<-EOF
    #!/bin/bash
    cd /home/ec2-user/
    docker-compose up -d
  EOF
  )
}

resource "aws_launch_template" "receptor_emergencia" {
  name     = "receptor_emergencia_launch_template"
  image_id = "ami-000807bf3875b8d6c"

  instance_market_options {
    market_type = "spot"
  }

  network_interfaces {
    security_groups             = [aws_security_group.ec2_http_access.id]
    associate_public_ip_address = true
  }

  instance_type = "m5.xlarge"
  user_data = base64encode(
    <<-EOF
    #!/bin/bash
    cd /home/ec2-user/
    docker-compose up -d
  EOF
  )
}

resource "aws_autoscaling_group" "normal_signals" {
  name                = "normal_signals"
  vpc_zone_identifier = [aws_subnet.az_1.id, aws_subnet.az_2.id]
  max_size            = 1
  min_size            = 1
  desired_capacity    = 1

  launch_template {
    id      = aws_launch_template.normal_receptor.id
    version = "$Latest"
  }

  instance_refresh {
    strategy = "Rolling"
    preferences {
      min_healthy_percentage = 0
    }
  }
  target_group_arns = [aws_lb_target_group.signals_receptor_1.arn, aws_lb_target_group.signals_receptor_2.arn, aws_lb_target_group.signals_receptor_3.arn, aws_lb_target_group.signals_receptor_4.arn]
}

resource "aws_autoscaling_group" "emergency_signals" {
  name                = "emergency_signals"
  vpc_zone_identifier = [aws_subnet.az_1.id, aws_subnet.az_2.id]
  max_size            = 0
  min_size            = 0
  desired_capacity    = 0

  launch_template {
    id      = aws_launch_template.receptor_emergencia.id
    version = "$Latest"
  }

  instance_refresh {
    strategy = "Rolling"
    preferences {
      min_healthy_percentage = 0
    }
  }
  target_group_arns = [aws_lb_target_group.emergency_receptor.arn]
}
