defmodule TaniaCore.Inventory.Material do
  use Ecto.Schema
  import Ecto.Changeset

  @material_types ~w(seed agrochemical growing_medium label_crop_support seeding_container post_harvest_supply plant other)

  @plant_types ~w(vegetable fruit herb flower tree ornament crop tuber)
  @chemical_types ~w(disinfectant fertilizer hormone manure pesticide)
  @container_types ~w(tray pot)

  @quantity_units_by_type %{
    "seed" => ~w(seeds packets gram kilogram),
    "agrochemical" => ~w(packets bottles bags),
    "growing_medium" => ~w(bags cubic_metre),
    "label_crop_support" => ~w(pieces),
    "seeding_container" => ~w(pieces),
    "post_harvest_supply" => ~w(pieces),
    "plant" => ~w(units packets),
    "other" => ~w(pieces)
  }

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "materials" do
    field :name, :string
    field :type, :string
    field :type_data, :map, default: %{}
    field :quantity_value, :float, default: 0.0
    field :quantity_unit, :string
    field :price_per_unit, :float
    field :price_currency, :string, default: "USD"
    field :expiration_date, :date
    field :notes, :string
    field :produced_by, :string

    belongs_to :farm, TaniaCore.Farming.Farm

    timestamps(type: :utc_datetime)
  end

  def changeset(material, attrs) do
    material
    |> cast(attrs, [
      :name,
      :type,
      :type_data,
      :quantity_value,
      :quantity_unit,
      :price_per_unit,
      :price_currency,
      :expiration_date,
      :notes,
      :produced_by,
      :farm_id
    ])
    |> validate_required([:name, :type, :farm_id])
    |> validate_inclusion(:type, @material_types)
    |> validate_number(:quantity_value, greater_than_or_equal_to: 0)
    |> validate_number(:price_per_unit, greater_than_or_equal_to: 0)
  end

  def material_types, do: @material_types
  def plant_types, do: @plant_types
  def chemical_types, do: @chemical_types
  def container_types, do: @container_types
  def quantity_units_for_type(type), do: Map.get(@quantity_units_by_type, type, ~w(pieces))
end
