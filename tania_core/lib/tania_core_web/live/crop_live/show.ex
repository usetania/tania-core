defmodule TaniaCoreWeb.CropLive.Show do
  use TaniaCoreWeb, :live_view

  alias TaniaCore.Farming
  alias TaniaCore.Growth

  @impl true
  def mount(%{"id" => id}, _session, socket) do
    crop = Growth.get_crop!(id)
    farm = Farming.get_user_farm(socket.assigns.current_scope.user.id)
    areas = if farm, do: Farming.list_areas(farm.id), else: []

    {:ok,
     socket
     |> assign(:page_title, "Crop: #{crop.batch_id}")
     |> assign(:active_page, :crops)
     |> assign(:crop, crop)
     |> assign(:farm, farm)
     |> assign(:areas, areas)
     |> assign(:active_modal, nil)
     |> assign(:action_form, to_form(%{}))}
  end

  # --- Action modals ---

  @impl true
  def handle_event("show_modal", %{"action" => action}, socket) do
    {:noreply, assign(socket, :active_modal, action)}
  end

  def handle_event("close_modal", _params, socket) do
    {:noreply, assign(socket, :active_modal, nil)}
  end

  # --- Move ---

  def handle_event("move_crop", params, socket) do
    crop = socket.assigns.crop
    user_id = socket.assigns.current_scope.user.id
    quantity = parse_int(params["quantity"])
    dst_area_id = params["dst_area_id"]
    src_area_id = params["src_area_id"]

    case Growth.move_crop(crop, quantity, dst_area_id, src_area_id, user_id) do
      {:ok, _movement} ->
        crop = Growth.get_crop!(crop.id)

        {:noreply,
         socket
         |> assign(:crop, crop)
         |> assign(:active_modal, nil)
         |> put_flash(:info, "Crop moved successfully")}

      {:error, _} ->
        {:noreply, put_flash(socket, :error, "Could not move crop")}
    end
  end

  # --- Harvest ---

  def handle_event("harvest_crop", params, socket) do
    crop = socket.assigns.crop
    user_id = socket.assigns.current_scope.user.id

    attrs = %{
      "harvest_type" => params["harvest_type"],
      "quantity" => params["quantity"],
      "produced_weight" => params["produced_weight"],
      "produced_unit" => params["produced_unit"],
      "area_id" => params["area_id"]
    }

    case Growth.harvest_crop(crop, attrs, user_id) do
      {:ok, _harvest} ->
        crop = Growth.get_crop!(crop.id)

        {:noreply,
         socket
         |> assign(:crop, crop)
         |> assign(:active_modal, nil)
         |> put_flash(:info, "Harvest recorded")}

      {:error, _} ->
        {:noreply, put_flash(socket, :error, "Could not record harvest")}
    end
  end

  # --- Dump ---

  def handle_event("dump_crop", params, socket) do
    crop = socket.assigns.crop
    user_id = socket.assigns.current_scope.user.id
    quantity = parse_int(params["quantity"])
    area_id = params["area_id"]

    case Growth.dump_crop(crop, quantity, area_id, user_id) do
      {:ok, _dump} ->
        crop = Growth.get_crop!(crop.id)

        {:noreply,
         socket
         |> assign(:crop, crop)
         |> assign(:active_modal, nil)
         |> put_flash(:info, "Dump recorded")}

      {:error, _} ->
        {:noreply, put_flash(socket, :error, "Could not record dump")}
    end
  end

  # --- Care actions ---

  def handle_event("water_crop", _params, socket) do
    crop = socket.assigns.crop
    user_id = socket.assigns.current_scope.user.id

    case Growth.water_crop(crop, nil, user_id) do
      {:ok, crop} ->
        crop = Growth.get_crop!(crop.id)
        {:noreply, socket |> assign(:crop, crop) |> put_flash(:info, "Watering recorded")}

      {:error, _} ->
        {:noreply, put_flash(socket, :error, "Could not record watering")}
    end
  end

  def handle_event("fertilize_crop", _params, socket) do
    crop = socket.assigns.crop
    user_id = socket.assigns.current_scope.user.id

    case Growth.fertilize_crop(crop, user_id) do
      {:ok, crop} ->
        crop = Growth.get_crop!(crop.id)
        {:noreply, socket |> assign(:crop, crop) |> put_flash(:info, "Fertilizing recorded")}

      {:error, _} ->
        {:noreply, put_flash(socket, :error, "Could not record fertilizing")}
    end
  end

  def handle_event("prune_crop", _params, socket) do
    crop = socket.assigns.crop
    user_id = socket.assigns.current_scope.user.id

    case Growth.prune_crop(crop, user_id) do
      {:ok, crop} ->
        crop = Growth.get_crop!(crop.id)
        {:noreply, socket |> assign(:crop, crop) |> put_flash(:info, "Pruning recorded")}

      {:error, _} ->
        {:noreply, put_flash(socket, :error, "Could not record pruning")}
    end
  end

  def handle_event("pesticide_crop", _params, socket) do
    crop = socket.assigns.crop
    user_id = socket.assigns.current_scope.user.id

    case Growth.pesticide_crop(crop, user_id) do
      {:ok, crop} ->
        crop = Growth.get_crop!(crop.id)

        {:noreply,
         socket |> assign(:crop, crop) |> put_flash(:info, "Pesticide application recorded")}

      {:error, _} ->
        {:noreply, put_flash(socket, :error, "Could not record pesticide application")}
    end
  end

  defp parse_int(nil), do: 0
  defp parse_int(""), do: 0
  defp parse_int(val) when is_binary(val), do: String.to_integer(val)
  defp parse_int(val) when is_integer(val), do: val

  defp format_datetime(nil), do: "-"
  defp format_datetime(dt), do: Calendar.strftime(dt, "%b %d, %Y %H:%M")

  defp activity_icon("create"), do: "hero-plus-circle"
  defp activity_icon("move"), do: "hero-arrow-right"
  defp activity_icon("harvest"), do: "hero-scissors"
  defp activity_icon("dump"), do: "hero-trash"
  defp activity_icon("water"), do: "hero-beaker"
  defp activity_icon("fertilize"), do: "hero-sparkles"
  defp activity_icon("prune"), do: "hero-scissors"
  defp activity_icon("pesticide"), do: "hero-shield-check"
  defp activity_icon("photo"), do: "hero-camera"
  defp activity_icon(_), do: "hero-clock"

  @impl true
  def render(assigns) do
    ~H"""
    <Layouts.app flash={@flash} current_scope={@current_scope} active_page={:crops}>
      <div class="mb-4">
        <.link navigate="/crops" class="text-sm text-primary hover:underline">
          <.icon name="hero-arrow-left" class="size-4" /> {gettext("Back to Crops")}
        </.link>
      </div>

      <.header>
        {gettext("Crop Batch")}: {@crop.batch_id}
        <:subtitle>
          <div class="flex gap-2 mt-1">
            <.badge variant={if(@crop.status == "active", do: "success", else: "default")}>
              {String.capitalize(@crop.status)}
            </.badge>
            <.badge variant={if(@crop.type == "seeding", do: "info", else: "success")}>
              {String.capitalize(@crop.type)}
            </.badge>
          </div>
        </:subtitle>
      </.header>

      <%!-- Info cards --%>
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 mt-6">
        <.stat_card label={gettext("Quantity")} value={@crop.quantity} icon="hero-hashtag" />
        <.stat_card
          label={gettext("Material")}
          value={if @crop.material, do: @crop.material.name, else: "-"}
          icon="hero-archive-box"
        />
        <.stat_card
          label={gettext("Last Watered")}
          value={format_datetime(@crop.last_watered)}
          icon="hero-beaker"
        />
        <.stat_card
          label={gettext("Last Fertilized")}
          value={format_datetime(@crop.last_fertilized)}
          icon="hero-sparkles"
        />
      </div>

      <%!-- Action buttons --%>
      <div :if={@crop.status == "active"} class="flex flex-wrap gap-2 mt-6">
        <.button phx-click={JS.push("show_modal", value: %{action: "move"})}>
          <.icon name="hero-arrow-right" class="size-4" /> {gettext("Move")}
        </.button>
        <.button phx-click={JS.push("show_modal", value: %{action: "harvest"})}>
          <.icon name="hero-scissors" class="size-4" /> {gettext("Harvest")}
        </.button>
        <.button phx-click={JS.push("show_modal", value: %{action: "dump"})} variant="danger">
          <.icon name="hero-trash" class="size-4" /> {gettext("Dump")}
        </.button>
        <.button phx-click="water_crop" variant="secondary">
          <.icon name="hero-beaker" class="size-4" /> {gettext("Water")}
        </.button>
        <.button phx-click="fertilize_crop" variant="secondary">
          <.icon name="hero-sparkles" class="size-4" /> {gettext("Fertilize")}
        </.button>
        <.button phx-click="prune_crop" variant="secondary">
          <.icon name="hero-scissors" class="size-4" /> {gettext("Prune")}
        </.button>
        <.button phx-click="pesticide_crop" variant="secondary">
          <.icon name="hero-shield-check" class="size-4" /> {gettext("Pesticide")}
        </.button>
      </div>

      <%!-- Activity timeline --%>
      <div class="mt-8">
        <h2 class="text-base font-semibold text-text mb-4">{gettext("Activity Timeline")}</h2>
        <.timeline :if={@crop.activities != []}>
          <.timeline_item
            :for={{activity, idx} <- Enum.with_index(@crop.activities)}
            icon={activity_icon(activity.activity_type)}
            time={format_datetime(activity.inserted_at)}
            last={idx == length(@crop.activities) - 1}
          >
            <span class="font-medium">{String.capitalize(activity.activity_type)}</span>
            <span :if={activity.description}>{" - #{activity.description}"}</span>
          </.timeline_item>
        </.timeline>
        <p :if={@crop.activities == []} class="text-sm text-text-light">
          {gettext("No activities recorded yet.")}
        </p>
      </div>

      <%!-- Move modal --%>
      <.modal
        :if={@active_modal == "move"}
        id="move-modal"
        show
        on_cancel={JS.push("close_modal")}
      >
        <h2 class="text-lg font-semibold mb-4">{gettext("Move Crop")}</h2>
        <form phx-submit="move_crop">
          <.input
            name="src_area_id"
            type="select"
            label={gettext("From Area")}
            options={Enum.map(@areas, &{&1.name, &1.id})}
            prompt={gettext("Select source area (optional)")}
            value=""
          />
          <.input
            name="dst_area_id"
            type="select"
            label={gettext("To Area")}
            options={Enum.map(@areas, &{&1.name, &1.id})}
            prompt={gettext("Select destination area")}
            value=""
          />
          <.input
            name="quantity"
            type="number"
            label={gettext("Quantity to move")}
            min="1"
            max={@crop.quantity}
            value={@crop.quantity}
          />
          <div class="mt-4 flex justify-end gap-2">
            <.button type="button" variant="secondary" phx-click="close_modal">
              {gettext("Cancel")}
            </.button>
            <.button type="submit">{gettext("Move")}</.button>
          </div>
        </form>
      </.modal>

      <%!-- Harvest modal --%>
      <.modal
        :if={@active_modal == "harvest"}
        id="harvest-modal"
        show
        on_cancel={JS.push("close_modal")}
      >
        <h2 class="text-lg font-semibold mb-4">{gettext("Harvest Crop")}</h2>
        <form phx-submit="harvest_crop">
          <.input
            name="harvest_type"
            type="select"
            label={gettext("Harvest Type")}
            options={[{gettext("All"), "all"}, {gettext("Partial"), "partial"}]}
            value="all"
          />
          <.input
            name="quantity"
            type="number"
            label={gettext("Quantity")}
            min="1"
            max={@crop.quantity}
            value={@crop.quantity}
          />
          <.input
            name="area_id"
            type="select"
            label={gettext("Area")}
            options={Enum.map(@areas, &{&1.name, &1.id})}
            prompt={gettext("Select area (optional)")}
            value=""
          />
          <.input
            name="produced_weight"
            type="number"
            label={gettext("Produced Weight")}
            step="any"
            value=""
          />
          <.input
            name="produced_unit"
            type="select"
            label={gettext("Weight Unit")}
            options={[{"Kg", "Kg"}, {"Gr", "Gr"}]}
            prompt={gettext("Select unit (optional)")}
            value=""
          />
          <div class="mt-4 flex justify-end gap-2">
            <.button type="button" variant="secondary" phx-click="close_modal">
              {gettext("Cancel")}
            </.button>
            <.button type="submit">{gettext("Harvest")}</.button>
          </div>
        </form>
      </.modal>

      <%!-- Dump modal --%>
      <.modal
        :if={@active_modal == "dump"}
        id="dump-modal"
        show
        on_cancel={JS.push("close_modal")}
      >
        <h2 class="text-lg font-semibold mb-4">{gettext("Dump Crop")}</h2>
        <form phx-submit="dump_crop">
          <.input
            name="quantity"
            type="number"
            label={gettext("Quantity to dump")}
            min="1"
            max={@crop.quantity}
            value={@crop.quantity}
          />
          <.input
            name="area_id"
            type="select"
            label={gettext("Area")}
            options={Enum.map(@areas, &{&1.name, &1.id})}
            prompt={gettext("Select area (optional)")}
            value=""
          />
          <div class="mt-4 flex justify-end gap-2">
            <.button type="button" variant="secondary" phx-click="close_modal">
              {gettext("Cancel")}
            </.button>
            <.button type="submit" variant="danger">{gettext("Dump")}</.button>
          </div>
        </form>
      </.modal>
    </Layouts.app>
    """
  end
end
