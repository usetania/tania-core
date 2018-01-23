<template lang="pug">
  .modal-mask
    .modal-wrapper
      .modal-container
        slot
          .modal-header
            span.h4.font-bold This is the Modal Header
          .modal-body
            span.h4.font-bold This is the Modal Body
</template>

<script>
export default {
  name: 'AppModalComponent'
}
</script>

