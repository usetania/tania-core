defmodule TaniaCore.Tasks.Task do
  use Ecto.Schema
  import Ecto.Changeset

  @priorities ~w(urgent normal)
  @statuses ~w(created completed cancelled)
  @categories ~w(general area reservoir crop finance inventory safety sanitation pestcontrol other)

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "tasks" do
    field :title, :string
    field :description, :string
    field :due_date, :utc_datetime
    field :completed_date, :utc_datetime
    field :cancelled_date, :utc_datetime
    field :priority, :string, default: "normal"
    field :status, :string, default: "created"
    field :category, :string, default: "general"
    field :is_due, :boolean, default: false
    field :domain, :string
    field :domain_data, :map, default: %{}
    field :asset_id, Ecto.UUID

    belongs_to :farm, TaniaCore.Farming.Farm

    timestamps(type: :utc_datetime)
  end

  def changeset(task, attrs) do
    task
    |> cast(attrs, [
      :title,
      :description,
      :due_date,
      :completed_date,
      :cancelled_date,
      :priority,
      :status,
      :category,
      :is_due,
      :domain,
      :domain_data,
      :asset_id,
      :farm_id
    ])
    |> validate_required([:title, :priority, :status, :category, :farm_id])
    |> validate_inclusion(:priority, @priorities)
    |> validate_inclusion(:status, @statuses)
    |> validate_inclusion(:category, @categories)
  end

  def priorities, do: @priorities
  def statuses, do: @statuses
  def categories, do: @categories
end
