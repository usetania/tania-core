defmodule TaniaCoreWeb.IntroLive.ReservoirLive do
  use TaniaCoreWeb, :live_view

  alias TaniaCore.Farming
  alias TaniaCore.Farming.Reservoir

  @impl true
  def mount(_params, _session, socket) do
    farm = Farming.get_user_farm(socket.assigns.current_scope.user.id)

    if farm do
      changeset = Farming.change_reservoir(%Reservoir{})

      {:ok,
       socket
       |> assign(:page_title, "Add a Reservoir")
       |> assign(:farm, farm)
       |> assign(:form, to_form(changeset))}
    else
      {:ok, push_navigate(socket, to: "/intro/farm")}
    end
  end

  @impl true
  def handle_event("validate", %{"reservoir" => params}, socket) do
    changeset =
      %Reservoir{}
      |> Farming.change_reservoir(params)
      |> Map.put(:action, :validate)

    {:noreply, assign(socket, :form, to_form(changeset))}
  end

  def handle_event("save", %{"reservoir" => params}, socket) do
    params = Map.put(params, "farm_id", socket.assigns.farm.id)
    user_id = socket.assigns.current_scope.user.id

    case Farming.create_reservoir(params, user_id) do
      {:ok, _reservoir} ->
        {:noreply, push_navigate(socket, to: "/intro/area")}

      {:error, changeset} ->
        {:noreply, assign(socket, :form, to_form(changeset))}
    end
  end

  def handle_event("skip", _params, socket) do
    {:noreply, push_navigate(socket, to: "/intro/area")}
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-background flex items-center justify-center p-4">
      <div class="w-full max-w-lg">
        <div class="text-center mb-8">
          <div class="flex justify-center mb-4">
            <div class="flex size-12 items-center justify-center rounded-xl bg-primary">
              <.icon name="hero-sun" class="size-7 text-white" />
            </div>
          </div>
          <h1 class="text-2xl font-bold text-text">{gettext("Add a Reservoir")}</h1>
          <p class="text-text-light mt-2">
            {gettext("Step 2 of 3. Add a water source for your farm.")}
          </p>
        </div>

        <%!-- Progress bar --%>
        <div class="flex gap-2 mb-8">
          <div class="h-1 flex-1 rounded-full bg-primary" />
          <div class="h-1 flex-1 rounded-full bg-primary" />
          <div class="h-1 flex-1 rounded-full bg-border" />
        </div>

        <div class="rounded-xl border border-border bg-white p-6 shadow-sm">
          <.form for={@form} phx-change="validate" phx-submit="save">
            <.input field={@form[:name]} label={gettext("Reservoir Name")} required />
            <.input
              field={@form[:water_source_type]}
              type="select"
              label={gettext("Water Source Type")}
              options={[
                {gettext("Tap"), "tap"},
                {gettext("Bucket"), "bucket"}
              ]}
              required
            />
            <.input
              field={@form[:capacity]}
              type="number"
              label={gettext("Capacity (Litres)")}
              step="any"
            />

            <div class="mt-6 flex gap-3">
              <.button type="button" variant="secondary" phx-click="skip" class="flex-1">
                {gettext("Skip")}
              </.button>
              <.button type="submit" class="flex-1">
                {gettext("Continue")} <.icon name="hero-arrow-right" class="size-4" />
              </.button>
            </div>
          </.form>
        </div>
      </div>
    </div>
    """
  end
end
