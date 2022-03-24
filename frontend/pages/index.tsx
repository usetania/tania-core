import type { NextPage } from "next";
import Link from "next/link";
import { Col, Form, Nav, Row, Table } from "react-bootstrap";

import Layout from "@/components/Layout";
import Panel from "@/components/Panel";
import TableTaskItem from "@/components/TableTaskItem";
import { dashboardData, cropsData, tasksData } from "@/data/index";

const Root: NextPage = () => {
  return (
    <Layout>
      <Row>
        <Col>
          <h3 className="pb-3">Dashboard</h3>
        </Col>
      </Row>
      <Row>
        <Panel
          title="At A Glance"
          md={6}
          lg={6}
        >
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
        </Panel>
        <Panel title="What is on Production" md={6} lg={6}>
          <>
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
          </>
        </Panel>
      </Row>
      <Row>
        <Panel title="Tasks">
          <>
            <div className="mb-3">
              <Link href="/tasks">See all Tasks</Link>
            </div>
            <Table responsive>
              <thead>
                <tr>
                  <th />
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
          </>
        </Panel>
      </Row>
    </Layout>
  );
};

export default Root;
