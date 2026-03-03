defmodule TaniaCore.Tasks do
  @moduledoc """
  The Tasks context manages farm operation tasks.
  """

  import Ecto.Query
  alias TaniaCore.Repo
  alias TaniaCore.Tasks.Task
  alias TaniaCore.Audit

  def list_tasks(farm_id, opts \\ []) do
    Task
    |> where([t], t.farm_id == ^farm_id)
    |> apply_filters(opts)
    |> order_by([t], desc: t.inserted_at)
    |> Repo.all()
  end

  def get_task!(id), do: Repo.get!(Task, id)

  def create_task(attrs, user_id \\ nil) do
    %Task{}
    |> Task.changeset(attrs)
    |> Repo.insert()
    |> tap_ok(fn task ->
      Audit.log_action("create", "task", task.id, attrs, user_id)
    end)
  end

  def update_task(%Task{} = task, attrs, user_id \\ nil) do
    task
    |> Task.changeset(attrs)
    |> Repo.update()
    |> tap_ok(fn task ->
      Audit.log_action("update", "task", task.id, attrs, user_id)
    end)
  end

  def delete_task(%Task{} = task, user_id \\ nil) do
    Repo.delete(task)
    |> tap_ok(fn _ ->
      Audit.log_action("delete", "task", task.id, %{}, user_id)
    end)
  end

  def complete_task(%Task{} = task, user_id \\ nil) do
    now = DateTime.utc_now() |> DateTime.truncate(:second)

    task
    |> Task.changeset(%{"status" => "completed", "completed_date" => now})
    |> Repo.update()
    |> tap_ok(fn task ->
      Audit.log_action("complete", "task", task.id, %{}, user_id)
    end)
  end

  def cancel_task(%Task{} = task, user_id \\ nil) do
    now = DateTime.utc_now() |> DateTime.truncate(:second)

    task
    |> Task.changeset(%{"status" => "cancelled", "cancelled_date" => now})
    |> Repo.update()
    |> tap_ok(fn task ->
      Audit.log_action("cancel", "task", task.id, %{}, user_id)
    end)
  end

  def check_due_tasks do
    now = DateTime.utc_now()

    from(t in Task,
      where: t.status == "created" and not is_nil(t.due_date) and t.due_date <= ^now,
      update: [set: [is_due: true]]
    )
    |> Repo.update_all([])
  end

  def change_task(%Task{} = task, attrs \\ %{}) do
    Task.changeset(task, attrs)
  end

  # --- Filters ---

  defp apply_filters(query, []), do: query

  defp apply_filters(query, [{:status, status} | rest]) when is_binary(status) do
    query
    |> where([t], t.status == ^status)
    |> apply_filters(rest)
  end

  defp apply_filters(query, [{:priority, priority} | rest]) when is_binary(priority) do
    query
    |> where([t], t.priority == ^priority)
    |> apply_filters(rest)
  end

  defp apply_filters(query, [{:category, category} | rest]) when is_binary(category) do
    query
    |> where([t], t.category == ^category)
    |> apply_filters(rest)
  end

  defp apply_filters(query, [{:is_due, true} | rest]) do
    query
    |> where([t], t.is_due == true)
    |> apply_filters(rest)
  end

  defp apply_filters(query, [{:due_today, true} | rest]) do
    today_start = Date.utc_today() |> DateTime.new!(~T[00:00:00], "Etc/UTC")
    today_end = Date.utc_today() |> DateTime.new!(~T[23:59:59], "Etc/UTC")

    query
    |> where([t], t.due_date >= ^today_start and t.due_date <= ^today_end)
    |> apply_filters(rest)
  end

  defp apply_filters(query, [{:due_this_week, true} | rest]) do
    today = Date.utc_today()
    week_start = Date.beginning_of_week(today)
    week_end = Date.end_of_week(today)
    start_dt = DateTime.new!(week_start, ~T[00:00:00], "Etc/UTC")
    end_dt = DateTime.new!(week_end, ~T[23:59:59], "Etc/UTC")

    query
    |> where([t], t.due_date >= ^start_dt and t.due_date <= ^end_dt)
    |> apply_filters(rest)
  end

  defp apply_filters(query, [{:due_this_month, true} | rest]) do
    today = Date.utc_today()
    month_start = Date.beginning_of_month(today)
    month_end = Date.end_of_month(today)
    start_dt = DateTime.new!(month_start, ~T[00:00:00], "Etc/UTC")
    end_dt = DateTime.new!(month_end, ~T[23:59:59], "Etc/UTC")

    query
    |> where([t], t.due_date >= ^start_dt and t.due_date <= ^end_dt)
    |> apply_filters(rest)
  end

  defp apply_filters(query, [_ | rest]), do: apply_filters(query, rest)

  defp tap_ok({:ok, record} = result, fun) do
    fun.(record)
    result
  end

  defp tap_ok(error, _fun), do: error
end
