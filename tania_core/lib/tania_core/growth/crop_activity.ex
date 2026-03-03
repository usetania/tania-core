defmodule TaniaCore.Growth.CropActivity do
  use Ecto.Schema
  import Ecto.Changeset

  @activity_types ~w(create move harvest dump water fertilize prune pesticide photo)

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "crop_activities" do
    field :activity_type, :string
    field :description, :string

    belongs_to :crop, TaniaCore.Growth.Crop

    timestamps(type: :utc_datetime, updated_at: false)
  end

  def changeset(activity, attrs) do
    activity
    |> cast(attrs, [:activity_type, :description, :crop_id])
    |> validate_required([:activity_type, :crop_id])
    |> validate_inclusion(:activity_type, @activity_types)
  end

  def activity_types, do: @activity_types
end
