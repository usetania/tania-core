defmodule TaniaCore.Inventory do
  @moduledoc """
  The Inventory context manages materials (seeds, chemicals, tools, etc).
  """

  import Ecto.Query
  alias TaniaCore.Repo
  alias TaniaCore.Inventory.Material
  alias TaniaCore.Audit

  def list_materials(farm_id) do
    Material
    |> where([m], m.farm_id == ^farm_id)
    |> order_by([m], asc: m.name)
    |> Repo.all()
  end

  def list_materials_by_type(farm_id, type) do
    Material
    |> where([m], m.farm_id == ^farm_id and m.type == ^type)
    |> order_by([m], asc: m.name)
    |> Repo.all()
  end

  def list_available_plant_types(farm_id) do
    Material
    |> where([m], m.farm_id == ^farm_id and m.type in ["seed", "plant"])
    |> order_by([m], asc: m.name)
    |> Repo.all()
  end

  def get_material!(id), do: Repo.get!(Material, id)

  def create_material(attrs, user_id \\ nil) do
    %Material{}
    |> Material.changeset(attrs)
    |> Repo.insert()
    |> tap_ok(fn material ->
      Audit.log_action("create", "material", material.id, attrs, user_id)
    end)
  end

  def update_material(%Material{} = material, attrs, user_id \\ nil) do
    material
    |> Material.changeset(attrs)
    |> Repo.update()
    |> tap_ok(fn material ->
      Audit.log_action("update", "material", material.id, attrs, user_id)
    end)
  end

  def delete_material(%Material{} = material, user_id \\ nil) do
    Repo.delete(material)
    |> tap_ok(fn _ ->
      Audit.log_action("delete", "material", material.id, %{}, user_id)
    end)
  end

  def change_material(%Material{} = material, attrs \\ %{}) do
    Material.changeset(material, attrs)
  end

  def quantity_units_for_type(type), do: Material.quantity_units_for_type(type)

  defp tap_ok({:ok, record} = result, fun) do
    fun.(record)
    result
  end

  defp tap_ok(error, _fun), do: error
end
