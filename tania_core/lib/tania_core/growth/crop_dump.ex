defmodule TaniaCore.Growth.CropDump do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "crop_dumps" do
    field :quantity, :integer

    belongs_to :crop, TaniaCore.Growth.Crop
    belongs_to :area, TaniaCore.Farming.Area

    timestamps(type: :utc_datetime, updated_at: false)
  end

  def changeset(dump, attrs) do
    dump
    |> cast(attrs, [:quantity, :crop_id, :area_id])
    |> validate_required([:quantity, :crop_id])
    |> validate_number(:quantity, greater_than: 0)
  end
end
