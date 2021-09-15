import {
  FaGripHorizontal,
  FaArchive,
  FaLeaf,
  FaHome,
  FaTint,
  FaClipboard,
  FaUser,
} from "react-icons/fa";

const dashboardData = [
  {
    name: "2 Areas",
    route: "/areas",
    icon: <FaGripHorizontal className="me-3" />,
  },
  {
    name: "1 Varieties",
    route: "/crops",
    icon: <FaLeaf className="me-3" />,
  },
  {
    name: "5 Tasks",
    route: "/materials",
    icon: <FaArchive className="me-3" />,
  },
];

const cropsData = [
  {
    id: 101,
    varieties: "Romaine",
    batchId: "rom-20mar",
    seedingDate: "20/03/2021",
    daysSinceSeeding: 177,
    qty: 126,
    qtyUnit: "Trays",
    lastWatering: "-",
  },
  {
    id: 102,
    varieties: "Romaine",
    batchId: "rom-24mar",
    seedingDate: "24/03/2021",
    daysSinceSeeding: 173,
    qty: 777,
    qtyUnit: "Pots",
    lastWatering: "-",
  },
];

const tasksData = [
  {
    id: 201,
    item: "Tempor commodo qui esse magna est culpa.",
    details:
      "Excepteur do non reprehenderit consequat eiusmod non nulla. Veniam esse nisi duis magna ex deserunt irure proident. Veniam mollit irure eu dolore quis esse cupidatat labore nulla. Nisi cillum veniam ipsum sunt ad aliqua sit labore. Adipisicing dolor quis eiusmod occaecat. Non velit commodo eu Lorem Lorem cillum deserunt sint dolore. Reprehenderit ipsum consequat minim voluptate fugiat incididunt proident id velit.",
    dueDate: "22/07/2021",
    priority: "normal",
    category: "sanitation",
  },
  {
    id: 202,
    item: "Laborum fugiat ea mollit aute.",
    details:
      "Consectetur non labore ut voluptate ullamco eu non nisi nisi dolor elit minim pariatur culpa. Fugiat ex ullamco sint dolore magna sit reprehenderit cillum veniam nisi dolor culpa ea. Aliquip veniam est sit eu in est occaecat nulla ipsum veniam.",
    dueDate: "22/07/2021",
    priority: "normal",
    category: "sanitation",
  },
  {
    id: 203,
    item: "Consequat incididunt sint do veniam velit eiusmod nisi ut aliqua dolor.",
    details:
      "Lorem aliqua nisi cillum nisi do consequat dolore minim Lorem ipsum. Fugiat tempor ipsum esse occaecat eiusmod. Fugiat sit pariatur aute voluptate enim aliquip magna. Id laborum labore deserunt exercitation mollit adipisicing do ipsum mollit. Amet proident pariatur in est cillum ex culpa anim mollit commodo cillum officia dolor ullamco.",
    dueDate: "20/07/2021",
    priority: "normal",
    category: "reservoir",
  },
  {
    id: 204,
    item: "Tempor excepteur qui in qui.",
    details:
      "Ut laborum cupidatat magna veniam ipsum exercitation. Quis do ipsum cillum ipsum fugiat ad nisi velit commodo consequat ipsum deserunt. Ad incididunt elit excepteur dolor laborum ex ipsum aliqua occaecat nisi nostrud qui sunt dolor. Labore ex qui non occaecat aliquip adipisicing est aliqua est non velit ullamco Lorem. Nulla deserunt sunt do culpa excepteur sunt tempor irure tempor enim voluptate. Ex et ullamco incididunt cupidatat sunt. Excepteur ullamco veniam adipisicing esse do reprehenderit ullamco.",
    dueDate: "20/07/2021",
    priority: "urgent",
    category: "safety",
  },
  {
    id: 205,
    item: "Veniam eiusmod dolor veniam excepteur fugiat anim eiusmod minim est proident consequat.",
    details:
      "Culpa non laborum laborum veniam excepteur. Minim velit deserunt sunt est. Ipsum irure aliqua nulla cupidatat officia in ad reprehenderit laborum deserunt fugiat esse ad pariatur.",
    dueDate: "30/03/2021",
    priority: "normal",
    category: "general",
  },
];

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

const reservoirData = [
  {
    id: 31,
    name: "River",
    createdOn: "20/03/2021",
    sourceType: "Tap/Well",
    capacity: 0,
    usedIn: "Organic Lettuce",
  },
  {
    id: 32,
    name: "Water Tank",
    createdOn: "19/07/2021",
    sourceType: "Water Tank/Cistern",
    capacity: 2,
    usedIn: "Organic Chilli",
  },
  {
    id: 33,
    name: "Something",
    createdOn: "07/09/2021",
    sourceType: "Tap/Well",
    capacity: 1,
    usedIn: "Organic Spinach",
  },
];

const notesData = [
  {
    id: 41,
    title: "Nulla eu aliquip veniam ea ea ad incididunt sint.",
    createdOn: "07/09/2021",
  },
  {
    id: 42,
    title:
      "Ut pariatur elit excepteur est in ad elit velit cillum mollit nostrud sint esse.",
    createdOn: "07/09/2021",
  },
  {
    id: 43,
    title: "Minim cillum irure non ad eu est voluptate ipsum voluptate.",
    createdOn: "19/07/2021",
  },
];

const areaData = [
  {
    id: 51,
    name: "Organic Lettuce",
    type: "Growing",
    size: 1,
    unit: "Ha",
    batches: 2,
    quantity: 779,
    edit: false,
  },
  {
    id: 52,
    name: "Organic Chilli",
    type: "Seeding",
    size: 2,
    unit: "Ha",
    batches: 0,
    quantity: 0,
    edit: true,
  },
];

const materialData = [
  {
    id: 61,
    category: "Seed",
    name: "Romaine",
    price: "2€",
    producedBy: "Kultiva",
    qty: 1000,
    additionalNotes: "",
  },
];

export {
  dashboardData,
  cropsData,
  tasksData,
  navData,
  reservoirData,
  notesData,
  areaData,
  materialData,
};
