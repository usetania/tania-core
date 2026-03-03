defmodule TaniaCore.Repo.Migrations.CreateMaterials do
  use Ecto.Migration

  def change do
    create table(:materials, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :name, :string, null: false
      add :type, :string, null: false
      add :type_data, :map, default: %{}
      add :quantity_value, :float, default: 0.0
      add :quantity_unit, :string
      add :price_per_unit, :float
      add :price_currency, :string, default: "USD"
      add :expiration_date, :date
      add :notes, :text
      add :produced_by, :string
      add :farm_id, references(:farms, type: :binary_id, on_delete: :delete_all), null: false

      timestamps(type: :utc_datetime)
    end

    create index(:materials, [:farm_id])
    create index(:materials, [:type])
  end
end
