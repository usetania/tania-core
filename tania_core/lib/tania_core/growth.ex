defmodule TaniaCore.Growth do
  @moduledoc """
  The Growth context manages crops and their lifecycle:
  create -> move -> harvest/dump -> archive.
  """

  import Ecto.Query
  alias TaniaCore.Repo
  alias TaniaCore.Growth.{Crop, CropMovement, CropHarvest, CropDump, CropActivity, CropPhoto}
  alias TaniaCore.Audit

  # --- Crops ---

  def list_crops(farm_id, status \\ "active") do
    Crop
    |> where([c], c.farm_id == ^farm_id and c.status == ^status)
    |> order_by([c], desc: c.inserted_at)
    |> preload([:material, :initial_area])
    |> Repo.all()
  end

  def get_crop!(id) do
    Crop
    |> Repo.get!(id)
    |> Repo.preload([
      :material,
      :initial_area,
      :movements,
      :harvests,
      :dumps,
      :photos,
      activities: {from(a in CropActivity, order_by: [desc: a.inserted_at]), []}
    ])
  end

  def create_crop(attrs, user_id \\ nil) do
    attrs = maybe_generate_batch_id(attrs)

    %Crop{}
    |> Crop.changeset(attrs)
    |> Repo.insert()
    |> tap_ok(fn crop ->
      log_activity(crop.id, "create", "Crop batch #{crop.batch_id} created")
      Audit.log_action("create", "crop", crop.id, attrs, user_id)
    end)
  end

  def update_crop(%Crop{} = crop, attrs, user_id \\ nil) do
    crop
    |> Crop.changeset(attrs)
    |> Repo.update()
    |> tap_ok(fn crop ->
      Audit.log_action("update", "crop", crop.id, attrs, user_id)
    end)
  end

  def delete_crop(%Crop{} = crop, user_id \\ nil) do
    Repo.delete(crop)
    |> tap_ok(fn _ ->
      Audit.log_action("delete", "crop", crop.id, %{}, user_id)
    end)
  end

  def change_crop(%Crop{} = crop, attrs \\ %{}) do
    Crop.changeset(crop, attrs)
  end

  # --- Move ---

  def move_crop(%Crop{} = crop, quantity, dst_area_id, src_area_id \\ nil, user_id \\ nil) do
    attrs = %{
      "quantity" => quantity,
      "crop_id" => crop.id,
      "dst_area_id" => dst_area_id,
      "src_area_id" => src_area_id
    }

    Repo.transaction(fn ->
      case %CropMovement{} |> CropMovement.changeset(attrs) |> Repo.insert() do
        {:ok, movement} ->
          log_activity(crop.id, "move", "Moved #{quantity} to new area")
          Audit.log_action("move", "crop", crop.id, attrs, user_id)
          movement

        {:error, changeset} ->
          Repo.rollback(changeset)
      end
    end)
  end

  # --- Harvest ---

  def harvest_crop(%Crop{} = crop, attrs, user_id \\ nil) do
    harvest_attrs = Map.put(attrs, "crop_id", crop.id)

    Repo.transaction(fn ->
      case %CropHarvest{} |> CropHarvest.changeset(harvest_attrs) |> Repo.insert() do
        {:ok, harvest} ->
          quantity = harvest.quantity
          new_quantity = max(crop.quantity - quantity, 0)

          crop
          |> Crop.changeset(%{"quantity" => new_quantity})
          |> maybe_archive(new_quantity)
          |> Repo.update!()

          desc = "Harvested #{quantity} (#{harvest.harvest_type})"
          log_activity(crop.id, "harvest", desc)
          Audit.log_action("harvest", "crop", crop.id, attrs, user_id)
          harvest

        {:error, changeset} ->
          Repo.rollback(changeset)
      end
    end)
  end

  # --- Dump ---

  def dump_crop(%Crop{} = crop, quantity, area_id \\ nil, user_id \\ nil) do
    attrs = %{
      "quantity" => quantity,
      "crop_id" => crop.id,
      "area_id" => area_id
    }

    Repo.transaction(fn ->
      case %CropDump{} |> CropDump.changeset(attrs) |> Repo.insert() do
        {:ok, dump} ->
          new_quantity = max(crop.quantity - quantity, 0)

          crop
          |> Crop.changeset(%{"quantity" => new_quantity})
          |> maybe_archive(new_quantity)
          |> Repo.update!()

          log_activity(crop.id, "dump", "Dumped #{quantity}")
          Audit.log_action("dump", "crop", crop.id, attrs, user_id)
          dump

        {:error, changeset} ->
          Repo.rollback(changeset)
      end
    end)
  end

  # --- Care activities ---

  def water_crop(%Crop{} = crop, area_description \\ nil, user_id \\ nil) do
    now = DateTime.utc_now() |> DateTime.truncate(:second)

    crop
    |> Crop.changeset(%{"last_watered" => now})
    |> Repo.update()
    |> tap_ok(fn _crop ->
      desc = if area_description, do: "Watered at #{area_description}", else: "Watered"
      log_activity(crop.id, "water", desc)
      Audit.log_action("water", "crop", crop.id, %{}, user_id)
    end)
  end

  def fertilize_crop(%Crop{} = crop, user_id \\ nil) do
    now = DateTime.utc_now() |> DateTime.truncate(:second)

    crop
    |> Crop.changeset(%{"last_fertilized" => now})
    |> Repo.update()
    |> tap_ok(fn _crop ->
      log_activity(crop.id, "fertilize", "Fertilized")
      Audit.log_action("fertilize", "crop", crop.id, %{}, user_id)
    end)
  end

  def prune_crop(%Crop{} = crop, user_id \\ nil) do
    now = DateTime.utc_now() |> DateTime.truncate(:second)

    crop
    |> Crop.changeset(%{"last_pruned" => now})
    |> Repo.update()
    |> tap_ok(fn _crop ->
      log_activity(crop.id, "prune", "Pruned")
      Audit.log_action("prune", "crop", crop.id, %{}, user_id)
    end)
  end

  def pesticide_crop(%Crop{} = crop, user_id \\ nil) do
    now = DateTime.utc_now() |> DateTime.truncate(:second)

    crop
    |> Crop.changeset(%{"last_pesticided" => now})
    |> Repo.update()
    |> tap_ok(fn _crop ->
      log_activity(crop.id, "pesticide", "Applied pesticide")
      Audit.log_action("pesticide", "crop", crop.id, %{}, user_id)
    end)
  end

  # --- Photos ---

  def add_photo(%Crop{} = crop, url, description \\ nil) do
    %CropPhoto{}
    |> CropPhoto.changeset(%{"url" => url, "description" => description, "crop_id" => crop.id})
    |> Repo.insert()
    |> tap_ok(fn _photo ->
      log_activity(crop.id, "photo", description || "Photo added")
    end)
  end

  def delete_photo(%CropPhoto{} = photo) do
    Repo.delete(photo)
  end

  def list_photos(crop_id) do
    CropPhoto
    |> where([p], p.crop_id == ^crop_id)
    |> order_by([p], desc: p.inserted_at)
    |> Repo.all()
  end

  # --- Activities ---

  def list_activities(crop_id) do
    CropActivity
    |> where([a], a.crop_id == ^crop_id)
    |> order_by([a], desc: a.inserted_at)
    |> Repo.all()
  end

  # --- Batch ID generation ---

  defp maybe_generate_batch_id(%{"batch_id" => id} = attrs) when is_binary(id) and id != "" do
    attrs
  end

  defp maybe_generate_batch_id(attrs) do
    date = Date.utc_today() |> Calendar.strftime("%y%m%d")
    suffix = :rand.uniform(999) |> Integer.to_string() |> String.pad_leading(3, "0")
    Map.put(attrs, "batch_id", "CROP-#{date}-#{suffix}")
  end

  # --- Helpers ---

  defp log_activity(crop_id, activity_type, description) do
    %CropActivity{}
    |> CropActivity.changeset(%{
      "crop_id" => crop_id,
      "activity_type" => activity_type,
      "description" => description
    })
    |> Repo.insert()
  end

  defp maybe_archive(changeset, 0) do
    Ecto.Changeset.put_change(changeset, :status, "archived")
  end

  defp maybe_archive(changeset, _quantity), do: changeset

  defp tap_ok({:ok, record} = result, fun) do
    fun.(record)
    result
  end

  defp tap_ok(error, _fun), do: error
end
