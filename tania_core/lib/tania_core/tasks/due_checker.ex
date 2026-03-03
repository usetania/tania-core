defmodule TaniaCore.Tasks.DueChecker do
  @moduledoc """
  Periodic job that marks tasks as overdue when their due date has passed.
  Runs via Quantum scheduler every 15 minutes.
  """

  def run do
    TaniaCore.Tasks.check_due_tasks()
  end
end
