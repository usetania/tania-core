<template>
  <div id="app">
    <div class="row">
      <div class="col-md-3 col-lg-2 d-none d-md-block">
        <AppAside />
      </div>
      <div class="col-xs-12 col-sm-12 col-md-9 col-lg-10">
        <div class="main-content">
          <div class="app-content-body">
            <router-view/>

          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import AppAside from '@/components/AppAside.vue';

export default {
  components: {
    AppAside,
  },
};
</script>

<style lang="scss">

</style>
