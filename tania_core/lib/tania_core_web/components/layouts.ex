defmodule TaniaCoreWeb.Layouts do
  @moduledoc """
  This module holds layouts and related functionality
  used by your application.
  """
  use TaniaCoreWeb, :html

  embed_templates "layouts/*"

  @doc """
  Renders the sidebar app layout with navigation.
  Used for authenticated pages.
  """
  attr :flash, :map, required: true, doc: "the map of flash messages"

  attr :current_scope, :map,
    default: nil,
    doc: "the current [scope](https://hexdocs.pm/phoenix/scopes.html)"

  attr :active_page, :atom, default: nil, doc: "the current active page for nav highlighting"

  slot :inner_block, required: true

  def app(assigns) do
    ~H"""
    <div class="min-h-screen bg-background">
      <!-- Mobile hamburger overlay -->
      <div
        id="mobile-sidebar-overlay"
        class="fixed inset-0 z-30 bg-black/50 lg:hidden hidden"
        phx-click={hide_sidebar()}
      />
      
    <!-- Sidebar -->
      <aside
        id="sidebar"
        class="fixed inset-y-0 left-0 z-40 flex w-64 flex-col bg-sidebar transition-transform duration-300 -translate-x-full lg:translate-x-0"
      >
        <!-- Logo -->
        <div class="flex h-16 items-center gap-3 px-6">
          <div class="flex size-8 items-center justify-center rounded-lg bg-white/10">
            <.icon name="hero-sun" class="size-5 text-white" />
          </div>
          <span class="text-lg font-bold text-white">Tania</span>
        </div>
        
    <!-- Navigation -->
        <nav class="flex-1 space-y-1 px-3 py-4">
          <.sidebar_link
            icon="hero-home"
            label={gettext("Dashboard")}
            href="/dashboard"
            active={@active_page == :dashboard}
          />
          <.sidebar_link
            icon="hero-map"
            label={gettext("Areas")}
            href="/areas"
            active={@active_page == :areas}
          />
          <.sidebar_link
            icon="hero-beaker"
            label={gettext("Reservoirs")}
            href="/reservoirs"
            active={@active_page == :reservoirs}
          />
          <.sidebar_link
            icon="hero-squares-2x2"
            label={gettext("Crops")}
            href="/crops"
            active={@active_page == :crops}
          />
          <.sidebar_link
            icon="hero-archive-box"
            label={gettext("Materials")}
            href="/materials"
            active={@active_page == :materials}
          />
          <.sidebar_link
            icon="hero-clipboard-document-list"
            label={gettext("Tasks")}
            href="/tasks"
            active={@active_page == :tasks}
          />
        </nav>
        
    <!-- User section -->
        <div :if={@current_scope} class="border-t border-white/10 px-3 py-4 space-y-1">
          <.sidebar_link
            icon="hero-cog-6-tooth"
            label={gettext("Settings")}
            href={~p"/users/settings"}
          />
          <.link
            href={~p"/users/log-out"}
            method="delete"
            class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-white/70 hover:bg-sidebar-hover hover:text-white transition-colors"
          >
            <.icon name="hero-arrow-right-on-rectangle" class="size-5" />
            <span>{gettext("Log out")}</span>
          </.link>
        </div>
      </aside>
      
    <!-- Main content -->
      <div class="lg:pl-64">
        <!-- Top header bar -->
        <header class="sticky top-0 z-20 flex h-16 items-center gap-4 border-b border-border bg-white px-4 sm:px-6 lg:px-8">
          <button
            class="lg:hidden text-text-light hover:text-text cursor-pointer"
            phx-click={show_sidebar()}
          >
            <.icon name="hero-bars-3" class="size-6" />
          </button>
          <div class="flex-1" />
          <span :if={@current_scope} class="text-sm text-text-light">
            {@current_scope.user.email}
          </span>
        </header>
        
    <!-- Page content -->
        <main class="p-4 sm:p-6 lg:p-8">
          {render_slot(@inner_block)}
        </main>
      </div>
    </div>

    <.flash_group flash={@flash} />
    """
  end

  defp show_sidebar(js \\ %JS{}) do
    js
    |> JS.show(to: "#mobile-sidebar-overlay")
    |> JS.remove_class("-translate-x-full", to: "#sidebar")
    |> JS.add_class("translate-x-0", to: "#sidebar")
  end

  defp hide_sidebar(js \\ %JS{}) do
    js
    |> JS.hide(to: "#mobile-sidebar-overlay")
    |> JS.add_class("-translate-x-full", to: "#sidebar")
    |> JS.remove_class("translate-x-0", to: "#sidebar")
  end

  @doc """
  Shows the flash group with standard titles and content.
  """
  attr :flash, :map, required: true, doc: "the map of flash messages"
  attr :id, :string, default: "flash-group", doc: "the optional id of flash container"

  def flash_group(assigns) do
    ~H"""
    <div id={@id} aria-live="polite">
      <.flash kind={:info} flash={@flash} />
      <.flash kind={:error} flash={@flash} />

      <.flash
        id="client-error"
        kind={:error}
        title={gettext("We can't find the internet")}
        phx-disconnected={show(".phx-client-error #client-error") |> JS.remove_attribute("hidden")}
        phx-connected={hide("#client-error") |> JS.set_attribute({"hidden", ""})}
        hidden
      >
        {gettext("Attempting to reconnect")}
        <.icon name="hero-arrow-path" class="ml-1 size-3 motion-safe:animate-spin" />
      </.flash>

      <.flash
        id="server-error"
        kind={:error}
        title={gettext("Something went wrong!")}
        phx-disconnected={show(".phx-server-error #server-error") |> JS.remove_attribute("hidden")}
        phx-connected={hide("#server-error") |> JS.set_attribute({"hidden", ""})}
        hidden
      >
        {gettext("Attempting to reconnect")}
        <.icon name="hero-arrow-path" class="ml-1 size-3 motion-safe:animate-spin" />
      </.flash>
    </div>
    """
  end
end
