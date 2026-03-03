defmodule TaniaCore.Audit.LogEntry do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "audit_logs" do
    field :action, :string
    field :entity_type, :string
    field :entity_id, :binary_id
    field :changes, :map, default: %{}

    belongs_to :performer, TaniaCore.Accounts.User, foreign_key: :performed_by

    timestamps(type: :utc_datetime, updated_at: false)
  end

  def changeset(log_entry, attrs) do
    log_entry
    |> cast(attrs, [:action, :entity_type, :entity_id, :changes, :performed_by])
    |> validate_required([:action, :entity_type, :entity_id])
  end
end
