import Image from "next/image";
import { Col, Nav } from "react-bootstrap";

import iNavData from "../types/iNavData";
import { navData } from "../data";

const Sidebar = (): JSX.Element => {
  return (
    <Col md={3} lg={2} className="bg-sidebar d-none d-md-block px-0 min-vh-100">
      <aside>
        <Nav defaultActiveKey="/dashboard" className="flex-column">
          <Nav.Link href="/">
            <div className="d-flex justify-content-center py-3 mb-3">
              <Image
                src={"/images/logo.png"}
                layout="fixed"
                width={100}
                height={33}
              />
            </div>
          </Nav.Link>
          {navData &&
            navData.map(({ name, route, icon }: iNavData) => (
              <Nav.Item key={name}>
                <Nav.Link href={route} className="text-light">
                  <div className="d-flex align-items-center">
                    {icon}
                    <span>{name}</span>
                  </div>
                </Nav.Link>
              </Nav.Item>
            ))}
        </Nav>
      </aside>
    </Col>
  );
};

export default Sidebar;
