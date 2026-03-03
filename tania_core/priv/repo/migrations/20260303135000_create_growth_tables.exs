defmodule TaniaCore.Repo.Migrations.CreateGrowthTables do
  use Ecto.Migration

  def change do
    create table(:crops, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :batch_id, :string, null: false
      add :status, :string, null: false, default: "active"
      add :type, :string, null: false
      add :container_type, :string
      add :container_cell, :integer
      add :quantity, :integer, null: false, default: 0
      add :last_watered, :utc_datetime
      add :last_fertilized, :utc_datetime
      add :last_pruned, :utc_datetime
      add :last_pesticided, :utc_datetime
      add :farm_id, references(:farms, type: :binary_id, on_delete: :delete_all), null: false
      add :material_id, references(:materials, type: :binary_id, on_delete: :nilify_all)
      add :initial_area_id, references(:areas, type: :binary_id, on_delete: :nilify_all)

      timestamps(type: :utc_datetime)
    end

    create index(:crops, [:farm_id])
    create index(:crops, [:status])
    create unique_index(:crops, [:batch_id])

    create table(:crop_movements, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :quantity, :integer, null: false
      add :crop_id, references(:crops, type: :binary_id, on_delete: :delete_all), null: false
      add :src_area_id, references(:areas, type: :binary_id, on_delete: :nilify_all)
      add :dst_area_id, references(:areas, type: :binary_id, on_delete: :nilify_all), null: false

      timestamps(type: :utc_datetime, updated_at: false)
    end

    create index(:crop_movements, [:crop_id])

    create table(:crop_harvests, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :harvest_type, :string, null: false
      add :quantity, :integer, null: false
      add :produced_weight, :float
      add :produced_unit, :string
      add :crop_id, references(:crops, type: :binary_id, on_delete: :delete_all), null: false
      add :area_id, references(:areas, type: :binary_id, on_delete: :nilify_all)

      timestamps(type: :utc_datetime, updated_at: false)
    end

    create index(:crop_harvests, [:crop_id])

    create table(:crop_dumps, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :quantity, :integer, null: false
      add :crop_id, references(:crops, type: :binary_id, on_delete: :delete_all), null: false
      add :area_id, references(:areas, type: :binary_id, on_delete: :nilify_all)

      timestamps(type: :utc_datetime, updated_at: false)
    end

    create index(:crop_dumps, [:crop_id])

    create table(:crop_activities, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :activity_type, :string, null: false
      add :description, :text
      add :crop_id, references(:crops, type: :binary_id, on_delete: :delete_all), null: false

      timestamps(type: :utc_datetime, updated_at: false)
    end

    create index(:crop_activities, [:crop_id])

    create table(:crop_photos, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :url, :string, null: false
      add :description, :string
      add :crop_id, references(:crops, type: :binary_id, on_delete: :delete_all), null: false

      timestamps(type: :utc_datetime, updated_at: false)
    end

    create index(:crop_photos, [:crop_id])
  end
end
