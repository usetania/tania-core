import Image from "next/image";
import Link from "next/link";
import React from "react";
import { Container, Col, Nav, Navbar, Row } from "react-bootstrap";
import {
  FaHome,
  FaTint,
  FaGripHorizontal,
  FaArchive,
  FaLeaf,
  FaClipboard,
  FaUser,
  FaPowerOff,
} from "react-icons/fa";

const navData = [
  {
    name: "Dashboard",
    route: "/dashboard",
    icon: <FaHome className="me-3" />,
  },
  {
    name: "Reservoir",
    route: "/reservoir",
    icon: <FaTint className="me-3" />,
  },
  {
    name: "Areas",
    route: "/areas",
    icon: <FaGripHorizontal className="me-3" />,
  },
  {
    name: "Materials",
    route: "/materials",
    icon: <FaArchive className="me-3" />,
  },
  {
    name: "Crops",
    route: "/crops",
    icon: <FaLeaf className="me-3" />,
  },
  {
    name: "Tasks",
    route: "/tasks",
    icon: <FaClipboard className="me-3" />,
  },
  {
    name: "Account",
    route: "/account",
    icon: <FaUser className="me-3" />,
  },
];

interface iLayout {
  children: React.ReactNode;
}

const Layout = ({ children }: iLayout) => {
  return (
    <Row className="mx-0">
      <Col md={3} lg={2} className="bg-primary d-none d-md-block px-0">
        <aside>
          <Nav defaultActiveKey="/" className="flex-column">
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
              navData.map(({ name, route, icon }) => (
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
      <Col sm={12} md={9} lg={10} className="bg-gray px-0">
        <Navbar collapseOnSelect className="bg-light px-3 py-2" expand="lg">
          <Navbar.Brand href="/">
            <div className="d-flex justify-content-center d-md-none">
              <Image
                src={"/images/logobig.png"}
                layout="fixed"
                width={100}
                height={33}
              />
            </div>
          </Navbar.Brand>
          <Navbar.Toggle />
          <Navbar.Collapse>
            <Nav className="d-none d-md-block">
              <Nav.Link>Demo Farm</Nav.Link>
            </Nav>
            <Nav className="d-md-none">
              {navData &&
                navData.map(({ name, route, icon }) => (
                  <Nav.Item key={`mobile-${name}`}>
                    <Nav.Link href={route}>
                      <div className="d-flex align-items-center">
                        {icon}
                        <span>{name}</span>
                      </div>
                    </Nav.Link>
                  </Nav.Item>
                ))}
            </Nav>
            <div className="dropdown-divider d-md-none" />
            <Nav className="ms-auto">
              <Nav.Link className="">
                <div className="d-flex align-items-center">
                  <FaPowerOff className="me-3" />
                  <Link href="/">
                    <span className="text-decoration-none">Sign Out</span>
                  </Link>
                </div>
              </Nav.Link>
            </Nav>
          </Navbar.Collapse>
        </Navbar>
        <Container fluid className="py-4">
          {children}
        </Container>
        <footer className="mb-3">
          <Container fluid>
            Tania 1.7.0. Made for the ♥ of plants &copy; 2019 Copyright.
          </Container>
        </footer>
      </Col>
    </Row>
  );
};

export default Layout;
