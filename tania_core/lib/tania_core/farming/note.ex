defmodule TaniaCore.Farming.Note do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "notes" do
    field :content, :string
    field :notable_type, :string
    field :notable_id, :binary_id

    belongs_to :creator, TaniaCore.Accounts.User, foreign_key: :created_by

    timestamps(type: :utc_datetime)
  end

  def changeset(note, attrs) do
    note
    |> cast(attrs, [:content, :notable_type, :notable_id, :created_by])
    |> validate_required([:content, :notable_type, :notable_id])
  end
end
