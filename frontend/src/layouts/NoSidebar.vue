<template>
  <div class="row">
    <div class="col-xs-12 col-sm-12 col-md-9 col-lg-10">
      <div class="main-content">
        <div class="app-content-body">
          <router-view/>
          <Footer />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Footer from '@/components/Footer.vue';

export default {
  components: {
    Footer,
  },
};
</script>
