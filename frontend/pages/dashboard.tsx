import type { NextPage } from "next";
import Link from "next/link";
import { Card, Col, Form, Nav, Row, Table } from "react-bootstrap";

import Layout from "../components/Layout";
import TableTaskItem from "../components/TableTaskItem";
import { dashboardData, cropsData, tasksData } from "../data";

const Dashboard: NextPage = () => {
  return (
    <Layout>
      <Row>
        <Col>
          <h3 className="pb-3">Dashboard</h3>
        </Col>
      </Row>
      <Row>
        <Col xs={12} sm={12} md={6}>
          <Card className="mb-3">
            <Card.Body>
              <Card.Title>
                <h5>At A Glance</h5>
              </Card.Title>
              <Nav className="flex-column">
                {dashboardData.map(({ name, route, icon }) => (
                  <Nav.Item key={name} className="mb-1">
                    <div className="d-flex align-items-center">
                      {icon}
                      <Link href={route}>{name}</Link>
                    </div>
                  </Nav.Item>
                ))}
              </Nav>
            </Card.Body>
            <Card.Footer>You are using Tania 1.7 right now</Card.Footer>
          </Card>
        </Col>
        <Col xs={12} sm={12} md={6}>
          <Card className="mb-3">
            <Card.Body>
              <Card.Title>
                <h5>What is on Production</h5>
              </Card.Title>
              <div className="mb-3">
                <Link href="/crops">See all Crops</Link>
              </div>
              <Table responsive>
                <thead>
                  <tr>
                    <th>Varieties</th>
                    <th>Qty</th>
                  </tr>
                </thead>
                <tbody>
                  {cropsData &&
                    cropsData.map(({ id, varieties, qty }) => (
                      <tr key={id}>
                        <td>
                          <Link href={`/crops/${id}`}>{varieties}</Link>
                        </td>
                        <td>{qty}</td>
                      </tr>
                    ))}
                </tbody>
              </Table>
            </Card.Body>
          </Card>
        </Col>
      </Row>
      <Row>
        <Col xs={12} sm={12} md={12}>
          <Card>
            <Card.Body>
              <Card.Title>
                <h5>Tasks</h5>
              </Card.Title>
              <div className="mb-3">
                <Link href="/tasks">See all Tasks</Link>
              </div>
              <Table responsive>
                <thead>
                  <tr>
                    <th></th>
                    <th className="w-75">Items</th>
                    <th>Category</th>
                  </tr>
                </thead>
                <tbody>
                  {tasksData &&
                    tasksData.map(
                      ({ id, item, details, dueDate, priority, category }) => (
                        <tr key={id}>
                          <td>
                            <Form>
                              <Form.Check type="checkbox" />
                            </Form>
                          </td>
                          <td>
                            <TableTaskItem
                              id={id}
                              item={item}
                              details={details}
                              dueDate={dueDate}
                              priority={priority}
                            />
                          </td>
                          <td>
                            <span className="text-uppercase">{category}</span>
                          </td>
                        </tr>
                      )
                    )}
                </tbody>
              </Table>
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </Layout>
  );
};

export default Dashboard;
