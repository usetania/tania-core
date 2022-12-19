import Image from "next/image";
import Link from "next/link";
import React from "react";
import { Container, Col, Nav, Navbar, Row } from "react-bootstrap";
import { FaPowerOff } from "react-icons/fa";

import Footer from "@/components/Footer";
import Sidebar from "./Sidebar";
import { navData } from "../data";

interface iLayout {
  children: React.ReactNode;
}

const Layout = ({ children }: iLayout) => {
  return (
    <Row className="mx-0">
      <Sidebar />
      <Col sm={12} md={9} lg={10} className="bg-gray px-0">
        <Navbar collapseOnSelect className="bg-light px-3 py-2" expand="lg">
          <Navbar.Brand href="/">
            <div className="d-flex justify-content-center d-md-none">
              <Image alt="Tania logo" src={"/img/logo.png"} width={100} height={33} />
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
              <Nav.Link>
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
          <div style={{ minHeight: "calc(100vh - 152px)" }}>{children}</div>
        </Container>
        <Footer />
      </Col>
    </Row>
  );
};

export default Layout;
