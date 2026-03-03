defmodule TaniaCore.Farming.Area do
  use Ecto.Schema
  import Ecto.Changeset

  @area_types ~w(seeding growing)
  @area_locations ~w(outdoor indoor)
  @size_units ~w(m2 Ha)

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "areas" do
    field :name, :string
    field :size_value, :float
    field :size_unit, :string
    field :type, :string
    field :location, :string
    field :photo_url, :string

    belongs_to :farm, TaniaCore.Farming.Farm
    belongs_to :reservoir, TaniaCore.Farming.Reservoir

    timestamps(type: :utc_datetime)
  end

  def changeset(area, attrs) do
    area
    |> cast(attrs, [
      :name,
      :size_value,
      :size_unit,
      :type,
      :location,
      :photo_url,
      :farm_id,
      :reservoir_id
    ])
    |> validate_required([:name, :type, :location, :farm_id])
    |> validate_inclusion(:type, @area_types)
    |> validate_inclusion(:location, @area_locations)
    |> validate_inclusion(:size_unit, @size_units)
    |> validate_number(:size_value, greater_than: 0)
  end

  def area_types, do: @area_types
  def area_locations, do: @area_locations
  def size_units, do: @size_units
end
