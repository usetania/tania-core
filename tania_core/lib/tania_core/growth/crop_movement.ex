defmodule TaniaCore.Growth.CropMovement do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "crop_movements" do
    field :quantity, :integer

    belongs_to :crop, TaniaCore.Growth.Crop
    belongs_to :src_area, TaniaCore.Farming.Area
    belongs_to :dst_area, TaniaCore.Farming.Area

    timestamps(type: :utc_datetime, updated_at: false)
  end

  def changeset(movement, attrs) do
    movement
    |> cast(attrs, [:quantity, :crop_id, :src_area_id, :dst_area_id])
    |> validate_required([:quantity, :crop_id, :dst_area_id])
    |> validate_number(:quantity, greater_than: 0)
  end
end
