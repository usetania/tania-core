<template>
  <div class="container-fluid bottom-space">
    <div class="row">
      <div class="col">
        <h3 class="title-page">Dashboard</h3>
      </div>
    </div>
    <div class="row">
      <div class="col-xs-12 col-sm-12 col-md-6">
        <b-card
          title="At A Glance"
          footer="You are using Tania 1.8 right now"
          class="card-ui"
        >
          <div class="row">
            <div class="col-xs-12 col-sm-12">
              <i class="fa fa-grip-horizontal"></i>
              <router-link :to="{ name: 'FarmAreas' }">0 Areas</router-link>
            </div>
            <div class="col-xs-12 col-sm-12">
              <i class="fa fa-leaf"></i>
              <router-link :to="{ name: 'FarmCrops' }">0 Varieties</router-link>
            </div>
            <div class="col-xs-12 col-sm-12">
              <i class="fa fa-clipboard"></i>
              <router-link :to="{ name: 'Tasks' }">0 Tasks</router-link>
            </div>
          </div>
        </b-card>
      </div>
      <!-- crops status -->
      <div class="col-xs-12 col-sm-12 col-md-6">
        <b-card
          title="What is On Production"
          class="card-ui"
        >
          <router-link :to="{ name: 'FarmCrops' }" class="see-all">
            See all Crops
          </router-link>
          <table class="table">
            <thead>
              <tr>
                <th>Varieties</th>
                <th>Qty</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>
                  <router-link :to="{ name: 'FarmCrop', params: { id: 1 } }">
                    Rose
                  </router-link>
                </td>
                <td>
                  12
                </td>
              </tr>
            </tbody>
          </table>
        </b-card>
      </div>
      <!-- end of crop status -->
    </div>

    <!-- Task List -->
    <div class="row">
      <div class="col-xs-12 col-sm-12 col-md-12">
        <b-card
          title="Tasks"
          class="card-ui"
        >
          <router-link :to="{ name: 'Task' }" class="see-all">
            See all Tasks
          </router-link>
        </b-card>
      </div>
    </div>
    <!-- End of task list -->
  </div>
</template>

<style lang="scss" scoped>
h3.title-page {
  margin: 20px 0 30px 0;
}

.card-ui {
  margin-bottom: 20px;

  i {
    width: 30px;
  }

  .see-all {
    display: block;
    margin-bottom: 15px;
  }
}

.bottom-space {
  padding-bottom: 60px;
}
</style>
