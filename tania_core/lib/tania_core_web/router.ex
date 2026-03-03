defmodule TaniaCoreWeb.Router do
  use TaniaCoreWeb, :router

  import TaniaCoreWeb.UserAuth

  pipeline :browser do
    plug :accepts, ["html"]
    plug :fetch_session
    plug :fetch_live_flash
    plug :put_root_layout, html: {TaniaCoreWeb.Layouts, :root}
    plug :protect_from_forgery
    plug :put_secure_browser_headers
    plug :fetch_current_scope_for_user
    plug TaniaCoreWeb.Plugs.SetLocale
  end

  pipeline :api do
    plug :accepts, ["json"]
  end

  # Public routes
  scope "/", TaniaCoreWeb do
    pipe_through :browser

    get "/", PageController, :home
  end

  # Enable LiveDashboard and Swoosh mailbox preview in development
  if Application.compile_env(:tania_core, :dev_routes) do
    import Phoenix.LiveDashboard.Router

    scope "/dev" do
      pipe_through :browser

      live_dashboard "/dashboard", metrics: TaniaCoreWeb.Telemetry
      forward "/mailbox", Plug.Swoosh.MailboxPreview
    end
  end

  ## Authentication routes

  scope "/", TaniaCoreWeb do
    pipe_through [:browser, :redirect_if_user_is_authenticated]

    get "/users/register", UserRegistrationController, :new
    post "/users/register", UserRegistrationController, :create
  end

  scope "/", TaniaCoreWeb do
    pipe_through [:browser, :require_authenticated_user]

    get "/users/settings", UserSettingsController, :edit
    put "/users/settings", UserSettingsController, :update
    get "/users/settings/confirm-email/:token", UserSettingsController, :confirm_email
  end

  scope "/", TaniaCoreWeb do
    pipe_through [:browser]

    get "/users/log-in", UserSessionController, :new
    get "/users/log-in/:token", UserSessionController, :confirm
    post "/users/log-in", UserSessionController, :create
    delete "/users/log-out", UserSessionController, :delete
  end

  ## Onboarding routes (authenticated, no farm check)

  scope "/", TaniaCoreWeb do
    pipe_through [:browser, :require_authenticated_user]

    live_session :onboarding,
      on_mount: [{TaniaCoreWeb.UserAuth, :ensure_authenticated}] do
      live "/intro/farm", IntroLive.FarmLive, :new
      live "/intro/reservoir", IntroLive.ReservoirLive, :new
      live "/intro/area", IntroLive.AreaLive, :new
    end
  end

  ## App routes (authenticated + farm required)

  scope "/", TaniaCoreWeb do
    pipe_through [:browser, :require_authenticated_user]

    live_session :app,
      on_mount: [
        {TaniaCoreWeb.UserAuth, :ensure_authenticated},
        {TaniaCoreWeb.Hooks.EnsureHasFarm, :default}
      ] do
      live "/dashboard", DashboardLive, :index

      live "/areas", AreaLive.Index, :index
      live "/areas/new", AreaLive.Index, :new
      live "/areas/:id/edit", AreaLive.Index, :edit

      live "/reservoirs", ReservoirLive.Index, :index

      live "/crops", CropLive.Index, :index
      live "/crops/:id", CropLive.Show, :show

      live "/materials", MaterialLive.Index, :index

      live "/tasks", TaskLive.Index, :index
    end
  end
end
