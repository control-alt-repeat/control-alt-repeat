variable "stage_name" {
  type        = string
  description = "The name of the stage. Must be either 'test' or 'live'"

  validation {
    condition     = contains(["test", "live"], var.test_variable)
    error_message = "Valid values for var: test_variable are (test, live)."
  } 
}
