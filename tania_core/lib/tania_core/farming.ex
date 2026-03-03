defmodule TaniaCore.Farming do
  @moduledoc """
  The Farming context manages farms, areas, reservoirs, and notes.
  """

  import Ecto.Query
  alias TaniaCore.Repo
  alias TaniaCore.Farming.{Farm, Reservoir, Area, Note}
  alias TaniaCore.Audit

  # --- Farms ---

  def list_farms(user_id) do
    Farm
    |> where([f], f.user_id == ^user_id)
    |> order_by([f], desc: f.inserted_at)
    |> Repo.all()
  end

  def get_farm!(id), do: Repo.get!(Farm, id)

  def get_user_farm(user_id) do
    Farm
    |> where([f], f.user_id == ^user_id)
    |> limit(1)
    |> Repo.one()
  end

  def create_farm(attrs, user_id) do
    %Farm{}
    |> Farm.changeset(Map.put(attrs, "user_id", user_id))
    |> Repo.insert()
    |> tap_ok(fn farm ->
      Audit.log_action("create", "farm", farm.id, attrs, user_id)
    end)
  end

  def update_farm(%Farm{} = farm, attrs, user_id \\ nil) do
    farm
    |> Farm.changeset(attrs)
    |> Repo.update()
    |> tap_ok(fn farm ->
      Audit.log_action("update", "farm", farm.id, attrs, user_id)
    end)
  end

  def change_farm(%Farm{} = farm, attrs \\ %{}) do
    Farm.changeset(farm, attrs)
  end

  # --- Reservoirs ---

  def list_reservoirs(farm_id) do
    Reservoir
    |> where([r], r.farm_id == ^farm_id)
    |> order_by([r], asc: r.name)
    |> Repo.all()
  end

  def get_reservoir!(id), do: Repo.get!(Reservoir, id)

  def create_reservoir(attrs, user_id \\ nil) do
    %Reservoir{}
    |> Reservoir.changeset(attrs)
    |> Repo.insert()
    |> tap_ok(fn reservoir ->
      Audit.log_action("create", "reservoir", reservoir.id, attrs, user_id)
    end)
  end

  def update_reservoir(%Reservoir{} = reservoir, attrs, user_id \\ nil) do
    reservoir
    |> Reservoir.changeset(attrs)
    |> Repo.update()
    |> tap_ok(fn reservoir ->
      Audit.log_action("update", "reservoir", reservoir.id, attrs, user_id)
    end)
  end

  def delete_reservoir(%Reservoir{} = reservoir, user_id \\ nil) do
    Repo.delete(reservoir)
    |> tap_ok(fn _ ->
      Audit.log_action("delete", "reservoir", reservoir.id, %{}, user_id)
    end)
  end

  def change_reservoir(%Reservoir{} = reservoir, attrs \\ %{}) do
    Reservoir.changeset(reservoir, attrs)
  end

  # --- Areas ---

  def list_areas(farm_id) do
    Area
    |> where([a], a.farm_id == ^farm_id)
    |> order_by([a], asc: a.name)
    |> preload(:reservoir)
    |> Repo.all()
  end

  def get_area!(id) do
    Area
    |> Repo.get!(id)
    |> Repo.preload(:reservoir)
  end

  def create_area(attrs, user_id \\ nil) do
    %Area{}
    |> Area.changeset(attrs)
    |> Repo.insert()
    |> tap_ok(fn area ->
      Audit.log_action("create", "area", area.id, attrs, user_id)
    end)
  end

  def update_area(%Area{} = area, attrs, user_id \\ nil) do
    area
    |> Area.changeset(attrs)
    |> Repo.update()
    |> tap_ok(fn area ->
      Audit.log_action("update", "area", area.id, attrs, user_id)
    end)
  end

  def delete_area(%Area{} = area, user_id \\ nil) do
    Repo.delete(area)
    |> tap_ok(fn _ ->
      Audit.log_action("delete", "area", area.id, %{}, user_id)
    end)
  end

  def change_area(%Area{} = area, attrs \\ %{}) do
    Area.changeset(area, attrs)
  end

  # --- Notes (polymorphic) ---

  def list_notes(notable_type, notable_id) do
    Note
    |> where([n], n.notable_type == ^notable_type and n.notable_id == ^notable_id)
    |> order_by([n], desc: n.inserted_at)
    |> Repo.all()
  end

  def create_note(notable_type, notable_id, attrs, user_id \\ nil) do
    attrs =
      Map.merge(attrs, %{
        "notable_type" => notable_type,
        "notable_id" => notable_id,
        "created_by" => user_id
      })

    %Note{}
    |> Note.changeset(attrs)
    |> Repo.insert()
  end

  def delete_note(%Note{} = note) do
    Repo.delete(note)
  end

  def change_note(%Note{} = note, attrs \\ %{}) do
    Note.changeset(note, attrs)
  end

  # --- Helpers ---

  defp tap_ok({:ok, record} = result, fun) do
    fun.(record)
    result
  end

  defp tap_ok(error, _fun), do: error
end
