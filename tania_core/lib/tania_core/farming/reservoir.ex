defmodule TaniaCore.Farming.Reservoir do
  use Ecto.Schema
  import Ecto.Changeset

  @water_source_types ~w(tap bucket)

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "reservoirs" do
    field :name, :string
    field :water_source_type, :string
    field :capacity, :float

    belongs_to :farm, TaniaCore.Farming.Farm
    has_many :areas, TaniaCore.Farming.Area

    timestamps(type: :utc_datetime)
  end

  def changeset(reservoir, attrs) do
    reservoir
    |> cast(attrs, [:name, :water_source_type, :capacity, :farm_id])
    |> validate_required([:name, :water_source_type, :farm_id])
    |> validate_inclusion(:water_source_type, @water_source_types)
    |> validate_bucket_capacity()
  end

  defp validate_bucket_capacity(changeset) do
    if get_field(changeset, :water_source_type) == "bucket" do
      validate_required(changeset, [:capacity])
      |> validate_number(:capacity, greater_than: 0)
    else
      changeset
    end
  end

  def water_source_types, do: @water_source_types
end
