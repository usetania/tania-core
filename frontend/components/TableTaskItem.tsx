import { useState } from "react";
import Link from "next/link";
import { Collapse } from "react-bootstrap";
import { FaChevronDown } from "react-icons/fa";

export interface iTableTaskItem {
  id: number;
  item: string;
  details: string;
  dueDate: string;
  priority: string;
  category?: string;
}

const TableTaskItem = ({
  id,
  item,
  details,
  dueDate,
  priority,
}: iTableTaskItem): JSX.Element => {
  const [showDetail, setShowDetail] = useState(false);
  return (
    <>
      <div className="mb-3">
        <Link href={`/tasks/${id}`}>{item}</Link>
        <br />
        <small className="d-flex align-items-center">
          <a
            onClick={() => setShowDetail(!showDetail)}
            aria-controls={`task-item-${id}`}
            className="lh-lg pe-auto text-black text-decoration-none"
            style={{ cursor: "pointer" }}
          >
            Read Details
            <FaChevronDown className="ms-1" />
          </a>
        </small>
        <Collapse in={showDetail}>
          <small id={`task-item-${id}`}>{details}</small>
        </Collapse>
      </div>
      <small className="text-muted">
        {`Due date: ${dueDate}`}
        <br />
        {"Priority: "}
        <span className="text-uppercase">{priority}</span>
      </small>
    </>
  );
};

export default TableTaskItem;
