defmodule TaniaCore.Repo.Migrations.CreateFarmingTables do
  use Ecto.Migration

  def change do
    create table(:farms, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :name, :string, null: false
      add :type, :string, null: false
      add :latitude, :string
      add :longitude, :string
      add :country, :string
      add :city, :string
      add :is_active, :boolean, default: true, null: false
      add :user_id, references(:users, type: :binary_id, on_delete: :delete_all), null: false

      timestamps(type: :utc_datetime)
    end

    create index(:farms, [:user_id])

    create table(:reservoirs, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :name, :string, null: false
      add :water_source_type, :string, null: false
      add :capacity, :float
      add :farm_id, references(:farms, type: :binary_id, on_delete: :delete_all), null: false

      timestamps(type: :utc_datetime)
    end

    create index(:reservoirs, [:farm_id])

    create table(:areas, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :name, :string, null: false
      add :size_value, :float
      add :size_unit, :string
      add :type, :string, null: false
      add :location, :string, null: false
      add :photo_url, :string
      add :farm_id, references(:farms, type: :binary_id, on_delete: :delete_all), null: false

      add :reservoir_id, references(:reservoirs, type: :binary_id, on_delete: :nilify_all)

      timestamps(type: :utc_datetime)
    end

    create index(:areas, [:farm_id])
    create index(:areas, [:reservoir_id])

    create table(:notes, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :content, :text, null: false
      add :notable_type, :string, null: false
      add :notable_id, :binary_id, null: false
      add :created_by, references(:users, type: :binary_id, on_delete: :nilify_all)

      timestamps(type: :utc_datetime)
    end

    create index(:notes, [:notable_type, :notable_id])
  end
end
