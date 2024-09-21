variable "stage_name" {
  type        = string
  description = "The name of the stage. Must be either 'test' or 'live'"

  validation {
    condition     = contains(["test", "live"], var.stage_name)
    error_message = "Valid values for var: stage_name are (test, live)."
  } 
}
