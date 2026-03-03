defmodule TaniaCoreWeb.PageControllerTest do
  use TaniaCoreWeb.ConnCase

  test "GET /", %{conn: conn} do
    conn = get(conn, ~p"/")
    assert html_response(conn, 200) =~ "Open-source farm management"
  end
end
