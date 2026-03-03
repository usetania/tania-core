defmodule TaniaCore.Growth.Crop do
  use Ecto.Schema
  import Ecto.Changeset

  @statuses ~w(active archived)
  @crop_types ~w(seeding growing)
  @container_types ~w(tray pot)

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "crops" do
    field :batch_id, :string
    field :status, :string, default: "active"
    field :type, :string
    field :container_type, :string
    field :container_cell, :integer
    field :quantity, :integer, default: 0
    field :last_watered, :utc_datetime
    field :last_fertilized, :utc_datetime
    field :last_pruned, :utc_datetime
    field :last_pesticided, :utc_datetime

    belongs_to :farm, TaniaCore.Farming.Farm
    belongs_to :material, TaniaCore.Inventory.Material
    belongs_to :initial_area, TaniaCore.Farming.Area

    has_many :movements, TaniaCore.Growth.CropMovement
    has_many :harvests, TaniaCore.Growth.CropHarvest
    has_many :dumps, TaniaCore.Growth.CropDump
    has_many :activities, TaniaCore.Growth.CropActivity
    has_many :photos, TaniaCore.Growth.CropPhoto

    timestamps(type: :utc_datetime)
  end

  def changeset(crop, attrs) do
    crop
    |> cast(attrs, [
      :batch_id,
      :status,
      :type,
      :container_type,
      :container_cell,
      :quantity,
      :last_watered,
      :last_fertilized,
      :last_pruned,
      :last_pesticided,
      :farm_id,
      :material_id,
      :initial_area_id
    ])
    |> validate_required([:batch_id, :type, :quantity, :farm_id])
    |> validate_inclusion(:status, @statuses)
    |> validate_inclusion(:type, @crop_types)
    |> validate_inclusion(:container_type, @container_types)
    |> validate_number(:quantity, greater_than_or_equal_to: 0)
    |> unique_constraint(:batch_id)
  end

  def statuses, do: @statuses
  def crop_types, do: @crop_types
  def container_types, do: @container_types
end
