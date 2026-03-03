defmodule TaniaCore.Audit do
  @moduledoc """
  Context for audit logging. Records actions performed on entities.
  """

  import Ecto.Query
  alias TaniaCore.Repo
  alias TaniaCore.Audit.LogEntry

  @doc """
  Logs an action performed on an entity.

  ## Examples

      iex> log_action("create", "farm", farm_id, %{name: "My Farm"}, user_id)
      {:ok, %LogEntry{}}
  """
  def log_action(action, entity_type, entity_id, changes \\ %{}, performed_by \\ nil) do
    %LogEntry{}
    |> LogEntry.changeset(%{
      action: action,
      entity_type: entity_type,
      entity_id: entity_id,
      changes: changes,
      performed_by: performed_by
    })
    |> Repo.insert()
  end

  @doc """
  Lists audit logs for a specific entity.
  """
  def list_logs_for_entity(entity_type, entity_id) do
    LogEntry
    |> where([l], l.entity_type == ^entity_type and l.entity_id == ^entity_id)
    |> order_by([l], desc: l.inserted_at)
    |> Repo.all()
  end
end
