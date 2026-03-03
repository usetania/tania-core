defmodule TaniaCore.Repo.Migrations.CreateTasks do
  use Ecto.Migration

  def change do
    create table(:tasks, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :title, :string, null: false
      add :description, :text
      add :due_date, :utc_datetime
      add :completed_date, :utc_datetime
      add :cancelled_date, :utc_datetime
      add :priority, :string, null: false, default: "normal"
      add :status, :string, null: false, default: "created"
      add :category, :string, null: false, default: "general"
      add :is_due, :boolean, default: false
      add :domain, :string
      add :domain_data, :map, default: %{}
      add :asset_id, :binary_id
      add :farm_id, references(:farms, type: :binary_id, on_delete: :delete_all), null: false

      timestamps(type: :utc_datetime)
    end

    create index(:tasks, [:farm_id])
    create index(:tasks, [:status])
    create index(:tasks, [:due_date])
  end
end
