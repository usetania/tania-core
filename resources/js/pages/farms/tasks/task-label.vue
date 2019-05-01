<template lang="pug">
  span
    span(v-if="type == 'CATEGORY'")
      span.label.label-pestcontrol(v-if="task.category == 'PESTCONTROL'")
        translate PEST CONTROL
      span.label.label-sanitation(v-if="task.category == 'SANITATION'")
        translate SANITATION
      span.label.label-area(v-if="task.category == 'AREA'")
        translate AREA
      span.label.label-reservoir(v-if="task.category == 'RESERVOIR'")
        translate RESERVOIR
      span.label.label-safety(v-if="task.category == 'SAFETY'")
        translate SAFETY
      span.label.label-finance(v-if="task.category == 'FINANCE'")
        translate FINANCE
      span.label.label-crop(v-if="task.category == 'CROP'")
        translate CROP
      span.label.label-water(v-if="task.category == 'WATER'")
        translate WATER
      span.label.label-inventory(v-if="task.category == 'INVENTORY'")
        translate INVENTORY
      span.label.label-general(v-if="task.category == 'GENERAL'")
        translate GENERAL
      span.label.label-nutrient(v-if="task.category == 'NUTRIENT'")
        translate NUTRIENT
    span(v-if="type == 'PRIORITY'")
      span.status.status-urgent(v-if="task.priority == 'URGENT'")
        translate URGENT
      span.status.status-normal(v-if="task.priority == 'NORMAL'")
        translate NORMAL
</template>

<script>
export default {
  name: 'TaskLabel',
  props: ['type', 'task'],
};
</script>
