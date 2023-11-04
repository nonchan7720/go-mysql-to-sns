package config

import "github.com/go-playground/validator/v10"

var validate = validator.New(validator.WithRequiredStructEnabled())

func init() {
	validate.RegisterStructValidation(fifoRequiredValidation, SNS{}, SQS{})
}

func fifoRequiredValidation(sl validator.StructLevel) {
	switch inf := sl.Current().Interface().(type) {
	case SNS:
		for _, t := range inf.Topics {
			if t.IsFIFO() && !t.Transform.IsOutbox() && t.MessageGroupIdTemplate == "" {
				sl.ReportError(t.MessageGroupIdTemplate, "MessageGroupIdTemplate", "messageGroupIdTemplate", "required_fifo_message_group_id", "")
			}
		}
	case SQS:
		for _, q := range inf.Queues {
			if q.IsFIFO() && !q.Transform.IsOutbox() && q.MessageGroupIdTemplate == "" {
				sl.ReportError(q.MessageGroupIdTemplate, "MessageGroupIdTemplate", "messageGroupIdTemplate", "required_fifo_message_group_id", "")
			}
		}
	}
}
