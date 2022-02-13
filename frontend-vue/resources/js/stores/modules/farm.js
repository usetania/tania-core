import FarmArea from './farms/area';
import FarmCrop from './farms/crop';
import FarmFarm from './farms/farm';
import FarmInventories from './farms/inventories';
import FarmReservoir from './farms/reservoir';
import Task from './farms/task';

const state = Object.assign({},
  FarmArea.state,
  FarmCrop.state,
  FarmFarm.state,
  FarmInventories.state,
  FarmReservoir.state,
  Task.state);

const getters = Object.assign({},
  FarmArea.getters,
  FarmCrop.getters,
  FarmFarm.getters,
  FarmInventories.getters,
  FarmReservoir.getters,
  Task.getters);

const actions = Object.assign({},
  FarmArea.actions,
  FarmCrop.actions,
  FarmFarm.actions,
  FarmInventories.actions,
  FarmReservoir.actions,
  Task.actions);

const mutations = Object.assign({},
  FarmArea.mutations,
  FarmCrop.mutations,
  FarmFarm.mutations,
  FarmInventories.mutations,
  FarmReservoir.mutations,
  Task.mutations);

export default {
  state, getters, actions, mutations,
};
