defmodule TaniaCoreWeb.PageController do
  use TaniaCoreWeb, :controller

  def home(conn, _params) do
    render(conn, :home)
  end
end
