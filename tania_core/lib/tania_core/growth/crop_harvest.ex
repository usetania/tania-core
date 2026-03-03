defmodule TaniaCore.Growth.CropHarvest do
  use Ecto.Schema
  import Ecto.Changeset

  @harvest_types ~w(all partial)
  @produced_units ~w(Kg Gr)

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "crop_harvests" do
    field :harvest_type, :string
    field :quantity, :integer
    field :produced_weight, :float
    field :produced_unit, :string

    belongs_to :crop, TaniaCore.Growth.Crop
    belongs_to :area, TaniaCore.Farming.Area

    timestamps(type: :utc_datetime, updated_at: false)
  end

  def changeset(harvest, attrs) do
    harvest
    |> cast(attrs, [
      :harvest_type,
      :quantity,
      :produced_weight,
      :produced_unit,
      :crop_id,
      :area_id
    ])
    |> validate_required([:harvest_type, :quantity, :crop_id])
    |> validate_inclusion(:harvest_type, @harvest_types)
    |> validate_inclusion(:produced_unit, @produced_units)
    |> validate_number(:quantity, greater_than: 0)
  end

  def harvest_types, do: @harvest_types
  def produced_units, do: @produced_units
end
