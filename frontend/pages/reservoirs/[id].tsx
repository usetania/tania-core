import { useRef, useState } from "react";
import type { NextPage } from "next";
import { useRouter } from "next/router";
import {
  Button,
  Col,
  Form,
  InputGroup,
  ListGroup,
  Row,
  Table,
} from "react-bootstrap";
import { FaPaperPlane, FaPlus, FaTrash } from "react-icons/fa";

import ButtonIcon from "../../components/ButtonIcon";
import Layout from "../../components/Layout";
import ModalContainer from "../../components/ModalContainer";
import Panel from "../../components/Panel";
import TableTaskItem from "../../components/TableTaskItem";
import useModal from "../../hooks/useModal";
import { tasksData, notesData } from "../../data";

const ReservoirDetail: NextPage = () => {
  const router = useRouter();
  const { id } = router.query;
  const { modalOpen, showModal, closeModal } = useModal();
  const [dueDate, setDueDate] = useState("");
  const [priority, setPriority] = useState("");
  const [title, setTitle] = useState("");
  const [desc, setDesc] = useState("");
  const [selectedCategory, setSelectedCategory] = useState("");
  const [isError, setIsError] = useState(false);
  const target = useRef(null);

  const addTaskReservoir = () => {
    if (!dueDate || !priority || !title) {
      setIsError(true);
    } else {
      setIsError(false);
      closeModal();
    }
  };

  return (
    <Layout>
      <Row>
        <Col xs={6}>
          <h3 className="pb-3">River</h3>
        </Col>
        <Col xs={6}>
          <div className="d-flex justify-content-end">
            <ButtonIcon
              label="Add Tasks"
              icon={<FaPlus className="me-2" />}
              onClick={showModal}
              variant="primary"
            />
          </div>
        </Col>
      </Row>
      <Row>
        <Panel title="Basic Info" md={6} lg={6}>
          <>
            <Row className="mb-2">
              <Col>
                <small>Source Type</small>
                <div>
                  <strong>Tap/Well</strong>
                </div>
              </Col>
              <Col>
                <small>Capacity</small>
                <div>
                  <strong>0</strong>
                </div>
              </Col>
            </Row>
            <Row>
              <Col>
                <small>Used In</small>
                <div>
                  <strong>Organic lettuce</strong>
                </div>
              </Col>
              <Col>
                <small>Created On</small>
                <div>
                  <strong>20/03/2021</strong>
                </div>
              </Col>
            </Row>
          </>
        </Panel>
        <Panel title="Notes" md={6} lg={6}>
          <>
            <InputGroup className="mb-3">
              <Form.Control type="text" placeholder="Create a note" />
              <Button variant="secondary">
                <div className="d-flex align-items-center">
                  <FaPaperPlane />
                </div>
              </Button>
            </InputGroup>
            <ListGroup>
              {notesData &&
                notesData.map(({ id, title, createdOn }) => (
                  <ListGroup.Item key={id}>
                    <div className="d-flex align-items-center justify-content-between py-1">
                      <div>
                        <div className="mb-1">{title}</div>
                        <small className="text-muted">{createdOn}</small>
                      </div>
                      <div>
                        <FaTrash />
                      </div>
                    </div>
                  </ListGroup.Item>
                ))}
            </ListGroup>
          </>
        </Panel>
      </Row>
      <Row>
        <Panel title="Tasks">
          <>
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
      <ModalContainer
        title="Reservoir: Add New Task on River"
        isShow={modalOpen}
        handleCloseModal={closeModal}
        handleSubmitModal={addTaskReservoir}
      >
        <>
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>Due Date</Form.Label>
              <InputGroup ref={target}>
                <Form.Control
                  type="date"
                  value={dueDate}
                  onChange={(e) => setDueDate(e.target.value)}
                />
              </InputGroup>
              {isError && (
                <Form.Text className="text-danger">
                  The due date field is required
                </Form.Text>
              )}
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Is this task urgent?</Form.Label>
              <Form.Check
                type="radio"
                label="Yes"
                name="priority"
                onChange={() => setPriority("urgent")}
              />
              <Form.Check
                type="radio"
                label="No"
                name="priority"
                onChange={() => setPriority("normal")}
              />
              {isError && (
                <Form.Text className="text-danger">
                  The priority field is required
                </Form.Text>
              )}
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Task Category</Form.Label>
              <Form.Select
                onChange={(e) => setSelectedCategory(e.target.value)}
              >
                <option>Please select category</option>
                <option value="1">Reservoir</option>
                <option value="2">Pest Control</option>
                <option value="3">Safety</option>
                <option value="4">Sanitation</option>
              </Form.Select>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Title</Form.Label>
              <Form.Control
                type="text"
                onChange={(e) => setTitle(e.target.value)}
                value={title}
              />
              {isError && (
                <Form.Text className="text-danger">
                  The title field is required
                </Form.Text>
              )}
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Description</Form.Label>
              <Form.Control
                as="textarea"
                onChange={(e) => setDesc(e.target.value)}
                value={desc}
                style={{ height: "120px" }}
              />
            </Form.Group>
          </Form>
        </>
      </ModalContainer>
    </Layout>
  );
};

export default ReservoirDetail;
