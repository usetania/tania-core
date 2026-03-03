defmodule TaniaCore.Farming.Farm do
  use Ecto.Schema
  import Ecto.Changeset

  @farm_types ~w(organic hydroponic aquaponic mushroom livestock fisheries permaculture)

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "farms" do
    field :name, :string
    field :type, :string
    field :latitude, :string
    field :longitude, :string
    field :country, :string
    field :city, :string
    field :is_active, :boolean, default: true

    belongs_to :user, TaniaCore.Accounts.User
    has_many :areas, TaniaCore.Farming.Area
    has_many :reservoirs, TaniaCore.Farming.Reservoir

    timestamps(type: :utc_datetime)
  end

  def changeset(farm, attrs) do
    farm
    |> cast(attrs, [:name, :type, :latitude, :longitude, :country, :city, :is_active, :user_id])
    |> validate_required([:name, :type, :user_id])
    |> validate_inclusion(:type, @farm_types)
  end

  def farm_types, do: @farm_types
end
