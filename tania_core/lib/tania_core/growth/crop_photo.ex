defmodule TaniaCore.Growth.CropPhoto do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "crop_photos" do
    field :url, :string
    field :description, :string

    belongs_to :crop, TaniaCore.Growth.Crop

    timestamps(type: :utc_datetime, updated_at: false)
  end

  def changeset(photo, attrs) do
    photo
    |> cast(attrs, [:url, :description, :crop_id])
    |> validate_required([:url, :crop_id])
  end
end
