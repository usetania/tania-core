defmodule TaniaCore.Scheduler do
  @moduledoc """
  Quantum scheduler for background jobs.
  Jobs are configured in config/config.exs.
  """
  use Quantum, otp_app: :tania_core
end
