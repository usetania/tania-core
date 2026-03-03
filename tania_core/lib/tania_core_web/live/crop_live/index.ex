defmodule TaniaCoreWeb.CropLive.Index do
  use TaniaCoreWeb, :live_view

  alias TaniaCore.Farming
  alias TaniaCore.Inventory
  alias TaniaCore.Growth
  alias TaniaCore.Growth.Crop

  @impl true
  def mount(_params, _session, socket) do
    farm = Farming.get_user_farm(socket.assigns.current_scope.user.id)

    if farm do
      crops = Growth.list_crops(farm.id, "active")
      areas = Farming.list_areas(farm.id)
      materials = Inventory.list_available_plant_types(farm.id)

      {:ok,
       socket
       |> assign(:page_title, "Crops")
       |> assign(:active_page, :crops)
       |> assign(:farm, farm)
       |> assign(:areas, areas)
       |> assign(:materials, materials)
       |> assign(:tab, "active")
       |> assign(:crops_count, length(crops))
       |> stream(:crops, crops)
       |> assign(:crop, nil)
       |> assign(:show_form, false)}
    else
      {:ok, push_navigate(socket, to: ~p"/")}
    end
  end

  @impl true
  def handle_event("switch_tab", %{"tab" => tab}, socket) do
    crops = Growth.list_crops(socket.assigns.farm.id, tab)

    {:noreply,
     socket
     |> assign(:tab, tab)
     |> assign(:crops_count, length(crops))
     |> stream(:crops, crops, reset: true)}
  end

  def handle_event("new_crop", _params, socket) do
    {:noreply,
     socket
     |> assign(:crop, %Crop{})
     |> assign(:show_form, true)
     |> assign(:form, to_form(Growth.change_crop(%Crop{})))}
  end

  def handle_event("close_form", _params, socket) do
    {:noreply, assign(socket, :show_form, false)}
  end

  def handle_event("validate", %{"crop" => crop_params}, socket) do
    crop = socket.assigns.crop || %Crop{}
    changeset = Growth.change_crop(crop, crop_params)
    {:noreply, assign(socket, :form, to_form(Map.put(changeset, :action, :validate)))}
  end

  def handle_event("save", %{"crop" => crop_params}, socket) do
    params =
      crop_params
      |> Map.put("farm_id", socket.assigns.farm.id)

    user_id = socket.assigns.current_scope.user.id

    case Growth.create_crop(params, user_id) do
      {:ok, crop} ->
        crop = Growth.get_crop!(crop.id)

        {:noreply,
         socket
         |> stream_insert(:crops, crop)
         |> update(:crops_count, &(&1 + 1))
         |> assign(:show_form, false)
         |> put_flash(:info, "Crop batch created")}

      {:error, changeset} ->
        {:noreply, assign(socket, :form, to_form(changeset))}
    end
  end

  def handle_event("delete", %{"id" => id}, socket) do
    crop = Growth.get_crop!(id)
    user_id = socket.assigns.current_scope.user.id

    case Growth.delete_crop(crop, user_id) do
      {:ok, _} ->
        {:noreply,
         socket
         |> stream_delete(:crops, crop)
         |> update(:crops_count, &(&1 - 1))
         |> put_flash(:info, "Crop deleted")}

      {:error, _} ->
        {:noreply, put_flash(socket, :error, "Could not delete crop")}
    end
  end

  defp format_date(nil), do: "-"
  defp format_date(dt), do: Calendar.strftime(dt, "%b %d, %Y")

  @impl true
  def render(assigns) do
    ~H"""
    <Layouts.app flash={@flash} current_scope={@current_scope} active_page={:crops}>
      <.header>
        {gettext("Crops")}
        <:actions>
          <.button phx-click="new_crop">
            <.icon name="hero-plus" class="size-4" /> {gettext("New Crop Batch")}
          </.button>
        </:actions>
      </.header>

      <.tabs
        tabs={[
          %{label: gettext("Active"), value: "active", active: @tab == "active"},
          %{label: gettext("Archived"), value: "archived", active: @tab == "archived"}
        ]}
        on_click="switch_tab"
      />

      <div class="mt-6">
        <.table
          :if={@crops_count > 0}
          id="crops"
          rows={@streams.crops}
          row_click={fn {_id, crop} -> JS.navigate("/crops/#{crop.id}") end}
        >
          <:col :let={{_id, c}} label={gettext("Batch ID")}>
            <span class="font-mono text-xs">{c.batch_id}</span>
          </:col>
          <:col :let={{_id, c}} label={gettext("Material")}>
            {if c.material, do: c.material.name, else: "-"}
          </:col>
          <:col :let={{_id, c}} label={gettext("Type")}>
            <.badge variant={if(c.type == "seeding", do: "info", else: "success")}>
              {String.capitalize(c.type)}
            </.badge>
          </:col>
          <:col :let={{_id, c}} label={gettext("Quantity")}>
            {c.quantity}
          </:col>
          <:col :let={{_id, c}} label={gettext("Status")}>
            <.badge variant={if(c.status == "active", do: "success", else: "default")}>
              {String.capitalize(c.status)}
            </.badge>
          </:col>
          <:col :let={{_id, c}} label={gettext("Created")}>
            {format_date(c.inserted_at)}
          </:col>
          <:action :let={{_id, c}}>
            <.button
              variant="ghost"
              phx-click={JS.push("delete", value: %{id: c.id})}
              data-confirm="Are you sure? This will delete the crop batch and all its history."
            >
              <.icon name="hero-trash" class="size-4" />
            </.button>
          </:action>
        </.table>

        <.empty_state
          :if={@crops_count == 0}
          icon="hero-squares-2x2"
          message={
            if @tab == "active",
              do: gettext("No active crops. Create a crop batch to get started."),
              else: gettext("No archived crops.")
          }
        >
          <:actions>
            <.button :if={@tab == "active"} phx-click="new_crop">
              <.icon name="hero-plus" class="size-4" /> {gettext("New Crop Batch")}
            </.button>
          </:actions>
        </.empty_state>
      </div>

      <.modal :if={@show_form} id="crop-form-modal" show on_cancel={JS.push("close_form")}>
        <h2 class="text-lg font-semibold mb-4">{gettext("New Crop Batch")}</h2>
        <.form for={@form} phx-change="validate" phx-submit="save">
          <.input
            field={@form[:type]}
            type="select"
            label={gettext("Crop Type")}
            options={[
              {gettext("Seeding"), "seeding"},
              {gettext("Growing"), "growing"}
            ]}
            required
          />
          <.input
            field={@form[:material_id]}
            type="select"
            label={gettext("Plant Material")}
            options={Enum.map(@materials, &{&1.name, &1.id})}
            prompt={gettext("Select material")}
          />
          <.input
            field={@form[:initial_area_id]}
            type="select"
            label={gettext("Initial Area")}
            options={Enum.map(@areas, &{&1.name, &1.id})}
            prompt={gettext("Select area")}
          />
          <.input
            field={@form[:container_type]}
            type="select"
            label={gettext("Container")}
            options={[
              {gettext("Tray"), "tray"},
              {gettext("Pot"), "pot"}
            ]}
            prompt={gettext("Select container (optional)")}
          />
          <.input
            field={@form[:container_cell]}
            type="number"
            label={gettext("Container Cells")}
            min="1"
          />
          <.input
            field={@form[:quantity]}
            type="number"
            label={gettext("Quantity")}
            min="1"
            required
          />
          <div class="mt-4 flex justify-end gap-2">
            <.button type="button" variant="secondary" phx-click="close_form">
              {gettext("Cancel")}
            </.button>
            <.button type="submit">{gettext("Create")}</.button>
          </div>
        </.form>
      </.modal>
    </Layouts.app>
    """
  end
end
