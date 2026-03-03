defmodule TaniaCoreWeb.Hooks.EnsureHasFarm do
  @moduledoc """
  LiveView on_mount hook that checks if the authenticated user has at least one farm.
  Redirects to the onboarding flow if no farm exists.
  """

  import Phoenix.LiveView

  alias TaniaCore.Farming

  def on_mount(:default, _params, _session, socket) do
    user_id = socket.assigns.current_scope.user.id
    farm = Farming.get_user_farm(user_id)

    if farm do
      {:cont, socket}
    else
      {:halt, push_navigate(socket, to: "/intro/farm")}
    end
  end
end
