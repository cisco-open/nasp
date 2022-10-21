// automatically generated by the FlatBuffers compiler, do not modify


// @generated

use core::mem;
use core::cmp::Ordering;

extern crate flatbuffers;
use self::flatbuffers::{EndianScalar, Follow};

#[allow(unused_imports, dead_code)]
pub mod wasm {

  use core::mem;
  use core::cmp::Ordering;

  extern crate flatbuffers;
  use self::flatbuffers::{EndianScalar, Follow};
#[allow(unused_imports, dead_code)]
pub mod common {

  use core::mem;
  use core::cmp::Ordering;

  extern crate flatbuffers;
  use self::flatbuffers::{EndianScalar, Follow};

pub enum KeyValOffset {}
#[derive(Copy, Clone, PartialEq)]

pub struct KeyVal<'a> {
  pub _tab: flatbuffers::Table<'a>,
}

impl<'a> flatbuffers::Follow<'a> for KeyVal<'a> {
  type Inner = KeyVal<'a>;
  #[inline]
  fn follow(buf: &'a [u8], loc: usize) -> Self::Inner {
    Self { _tab: flatbuffers::Table { buf, loc } }
  }
}

impl<'a> KeyVal<'a> {
  pub const VT_KEY: flatbuffers::VOffsetT = 4;
  pub const VT_VALUE: flatbuffers::VOffsetT = 6;

  #[inline]
  pub fn init_from_table(table: flatbuffers::Table<'a>) -> Self {
    KeyVal { _tab: table }
  }
  #[allow(unused_mut)]
  pub fn create<'bldr: 'args, 'args: 'mut_bldr, 'mut_bldr>(
    _fbb: &'mut_bldr mut flatbuffers::FlatBufferBuilder<'bldr>,
    args: &'args KeyValArgs<'args>
  ) -> flatbuffers::WIPOffset<KeyVal<'bldr>> {
    let mut builder = KeyValBuilder::new(_fbb);
    if let Some(x) = args.value { builder.add_value(x); }
    if let Some(x) = args.key { builder.add_key(x); }
    builder.finish()
  }


  #[inline]
  pub fn key(&self) -> &'a str {
    self._tab.get::<flatbuffers::ForwardsUOffset<&str>>(KeyVal::VT_KEY, None).unwrap()
  }
  #[inline]
  pub fn key_compare_less_than(&self, o: &KeyVal) -> bool {
    self.key() < o.key()
  }

  #[inline]
  pub fn key_compare_with_value(&self, val: & str) -> ::core::cmp::Ordering {
    let key = self.key();
    key.cmp(val)
  }
  #[inline]
  pub fn value(&self) -> Option<&'a str> {
    self._tab.get::<flatbuffers::ForwardsUOffset<&str>>(KeyVal::VT_VALUE, None)
  }
}

impl flatbuffers::Verifiable for KeyVal<'_> {
  #[inline]
  fn run_verifier(
    v: &mut flatbuffers::Verifier, pos: usize
  ) -> Result<(), flatbuffers::InvalidFlatbuffer> {
    use self::flatbuffers::Verifiable;
    v.visit_table(pos)?
     .visit_field::<flatbuffers::ForwardsUOffset<&str>>("key", Self::VT_KEY, true)?
     .visit_field::<flatbuffers::ForwardsUOffset<&str>>("value", Self::VT_VALUE, false)?
     .finish();
    Ok(())
  }
}
pub struct KeyValArgs<'a> {
    pub key: Option<flatbuffers::WIPOffset<&'a str>>,
    pub value: Option<flatbuffers::WIPOffset<&'a str>>,
}
impl<'a> Default for KeyValArgs<'a> {
  #[inline]
  fn default() -> Self {
    KeyValArgs {
      key: None, // required field
      value: None,
    }
  }
}

pub struct KeyValBuilder<'a: 'b, 'b> {
  fbb_: &'b mut flatbuffers::FlatBufferBuilder<'a>,
  start_: flatbuffers::WIPOffset<flatbuffers::TableUnfinishedWIPOffset>,
}
impl<'a: 'b, 'b> KeyValBuilder<'a, 'b> {
  #[inline]
  pub fn add_key(&mut self, key: flatbuffers::WIPOffset<&'b  str>) {
    self.fbb_.push_slot_always::<flatbuffers::WIPOffset<_>>(KeyVal::VT_KEY, key);
  }
  #[inline]
  pub fn add_value(&mut self, value: flatbuffers::WIPOffset<&'b  str>) {
    self.fbb_.push_slot_always::<flatbuffers::WIPOffset<_>>(KeyVal::VT_VALUE, value);
  }
  #[inline]
  pub fn new(_fbb: &'b mut flatbuffers::FlatBufferBuilder<'a>) -> KeyValBuilder<'a, 'b> {
    let start = _fbb.start_table();
    KeyValBuilder {
      fbb_: _fbb,
      start_: start,
    }
  }
  #[inline]
  pub fn finish(self) -> flatbuffers::WIPOffset<KeyVal<'a>> {
    let o = self.fbb_.end_table(self.start_);
    self.fbb_.required(o, KeyVal::VT_KEY,"key");
    flatbuffers::WIPOffset::new(o.value())
  }
}

impl core::fmt::Debug for KeyVal<'_> {
  fn fmt(&self, f: &mut core::fmt::Formatter<'_>) -> core::fmt::Result {
    let mut ds = f.debug_struct("KeyVal");
      ds.field("key", &self.key());
      ds.field("value", &self.value());
      ds.finish()
  }
}
pub enum FlatNodeOffset {}
#[derive(Copy, Clone, PartialEq)]

pub struct FlatNode<'a> {
  pub _tab: flatbuffers::Table<'a>,
}

impl<'a> flatbuffers::Follow<'a> for FlatNode<'a> {
  type Inner = FlatNode<'a>;
  #[inline]
  fn follow(buf: &'a [u8], loc: usize) -> Self::Inner {
    Self { _tab: flatbuffers::Table { buf, loc } }
  }
}

impl<'a> FlatNode<'a> {
  pub const VT_NAME: flatbuffers::VOffsetT = 4;
  pub const VT_NAMESPACE: flatbuffers::VOffsetT = 6;
  pub const VT_LABELS: flatbuffers::VOffsetT = 8;
  pub const VT_OWNER: flatbuffers::VOffsetT = 10;
  pub const VT_WORKLOAD_NAME: flatbuffers::VOffsetT = 12;
  pub const VT_PLATFORM_METADATA: flatbuffers::VOffsetT = 14;
  pub const VT_ISTIO_VERSION: flatbuffers::VOffsetT = 16;
  pub const VT_MESH_ID: flatbuffers::VOffsetT = 18;
  pub const VT_APP_CONTAINERS: flatbuffers::VOffsetT = 20;
  pub const VT_CLUSTER_ID: flatbuffers::VOffsetT = 22;
  pub const VT_INSTANCE_IPS: flatbuffers::VOffsetT = 24;

  #[inline]
  pub fn init_from_table(table: flatbuffers::Table<'a>) -> Self {
    FlatNode { _tab: table }
  }
  #[allow(unused_mut)]
  pub fn create<'bldr: 'args, 'args: 'mut_bldr, 'mut_bldr>(
    _fbb: &'mut_bldr mut flatbuffers::FlatBufferBuilder<'bldr>,
    args: &'args FlatNodeArgs<'args>
  ) -> flatbuffers::WIPOffset<FlatNode<'bldr>> {
    let mut builder = FlatNodeBuilder::new(_fbb);
    if let Some(x) = args.instance_ips { builder.add_instance_ips(x); }
    if let Some(x) = args.cluster_id { builder.add_cluster_id(x); }
    if let Some(x) = args.app_containers { builder.add_app_containers(x); }
    if let Some(x) = args.mesh_id { builder.add_mesh_id(x); }
    if let Some(x) = args.istio_version { builder.add_istio_version(x); }
    if let Some(x) = args.platform_metadata { builder.add_platform_metadata(x); }
    if let Some(x) = args.workload_name { builder.add_workload_name(x); }
    if let Some(x) = args.owner { builder.add_owner(x); }
    if let Some(x) = args.labels { builder.add_labels(x); }
    if let Some(x) = args.namespace { builder.add_namespace(x); }
    if let Some(x) = args.name { builder.add_name(x); }
    builder.finish()
  }


  #[inline]
  pub fn name(&self) -> Option<&'a str> {
    self._tab.get::<flatbuffers::ForwardsUOffset<&str>>(FlatNode::VT_NAME, None)
  }
  #[inline]
  pub fn namespace(&self) -> Option<&'a str> {
    self._tab.get::<flatbuffers::ForwardsUOffset<&str>>(FlatNode::VT_NAMESPACE, None)
  }
  #[inline]
  pub fn labels(&self) -> Option<flatbuffers::Vector<'a, flatbuffers::ForwardsUOffset<KeyVal<'a>>>> {
    self._tab.get::<flatbuffers::ForwardsUOffset<flatbuffers::Vector<'a, flatbuffers::ForwardsUOffset<KeyVal>>>>(FlatNode::VT_LABELS, None)
  }
  #[inline]
  pub fn owner(&self) -> Option<&'a str> {
    self._tab.get::<flatbuffers::ForwardsUOffset<&str>>(FlatNode::VT_OWNER, None)
  }
  #[inline]
  pub fn workload_name(&self) -> Option<&'a str> {
    self._tab.get::<flatbuffers::ForwardsUOffset<&str>>(FlatNode::VT_WORKLOAD_NAME, None)
  }
  #[inline]
  pub fn platform_metadata(&self) -> Option<flatbuffers::Vector<'a, flatbuffers::ForwardsUOffset<KeyVal<'a>>>> {
    self._tab.get::<flatbuffers::ForwardsUOffset<flatbuffers::Vector<'a, flatbuffers::ForwardsUOffset<KeyVal>>>>(FlatNode::VT_PLATFORM_METADATA, None)
  }
  #[inline]
  pub fn istio_version(&self) -> Option<&'a str> {
    self._tab.get::<flatbuffers::ForwardsUOffset<&str>>(FlatNode::VT_ISTIO_VERSION, None)
  }
  #[inline]
  pub fn mesh_id(&self) -> Option<&'a str> {
    self._tab.get::<flatbuffers::ForwardsUOffset<&str>>(FlatNode::VT_MESH_ID, None)
  }
  #[inline]
  pub fn app_containers(&self) -> Option<flatbuffers::Vector<'a, flatbuffers::ForwardsUOffset<&'a str>>> {
    self._tab.get::<flatbuffers::ForwardsUOffset<flatbuffers::Vector<'a, flatbuffers::ForwardsUOffset<&'a str>>>>(FlatNode::VT_APP_CONTAINERS, None)
  }
  #[inline]
  pub fn cluster_id(&self) -> Option<&'a str> {
    self._tab.get::<flatbuffers::ForwardsUOffset<&str>>(FlatNode::VT_CLUSTER_ID, None)
  }
  #[inline]
  pub fn instance_ips(&self) -> Option<flatbuffers::Vector<'a, flatbuffers::ForwardsUOffset<&'a str>>> {
    self._tab.get::<flatbuffers::ForwardsUOffset<flatbuffers::Vector<'a, flatbuffers::ForwardsUOffset<&'a str>>>>(FlatNode::VT_INSTANCE_IPS, None)
  }
}

impl flatbuffers::Verifiable for FlatNode<'_> {
  #[inline]
  fn run_verifier(
    v: &mut flatbuffers::Verifier, pos: usize
  ) -> Result<(), flatbuffers::InvalidFlatbuffer> {
    use self::flatbuffers::Verifiable;
    v.visit_table(pos)?
     .visit_field::<flatbuffers::ForwardsUOffset<&str>>("name", Self::VT_NAME, false)?
     .visit_field::<flatbuffers::ForwardsUOffset<&str>>("namespace", Self::VT_NAMESPACE, false)?
     .visit_field::<flatbuffers::ForwardsUOffset<flatbuffers::Vector<'_, flatbuffers::ForwardsUOffset<KeyVal>>>>("labels", Self::VT_LABELS, false)?
     .visit_field::<flatbuffers::ForwardsUOffset<&str>>("owner", Self::VT_OWNER, false)?
     .visit_field::<flatbuffers::ForwardsUOffset<&str>>("workload_name", Self::VT_WORKLOAD_NAME, false)?
     .visit_field::<flatbuffers::ForwardsUOffset<flatbuffers::Vector<'_, flatbuffers::ForwardsUOffset<KeyVal>>>>("platform_metadata", Self::VT_PLATFORM_METADATA, false)?
     .visit_field::<flatbuffers::ForwardsUOffset<&str>>("istio_version", Self::VT_ISTIO_VERSION, false)?
     .visit_field::<flatbuffers::ForwardsUOffset<&str>>("mesh_id", Self::VT_MESH_ID, false)?
     .visit_field::<flatbuffers::ForwardsUOffset<flatbuffers::Vector<'_, flatbuffers::ForwardsUOffset<&'_ str>>>>("app_containers", Self::VT_APP_CONTAINERS, false)?
     .visit_field::<flatbuffers::ForwardsUOffset<&str>>("cluster_id", Self::VT_CLUSTER_ID, false)?
     .visit_field::<flatbuffers::ForwardsUOffset<flatbuffers::Vector<'_, flatbuffers::ForwardsUOffset<&'_ str>>>>("instance_ips", Self::VT_INSTANCE_IPS, false)?
     .finish();
    Ok(())
  }
}
pub struct FlatNodeArgs<'a> {
    pub name: Option<flatbuffers::WIPOffset<&'a str>>,
    pub namespace: Option<flatbuffers::WIPOffset<&'a str>>,
    pub labels: Option<flatbuffers::WIPOffset<flatbuffers::Vector<'a, flatbuffers::ForwardsUOffset<KeyVal<'a>>>>>,
    pub owner: Option<flatbuffers::WIPOffset<&'a str>>,
    pub workload_name: Option<flatbuffers::WIPOffset<&'a str>>,
    pub platform_metadata: Option<flatbuffers::WIPOffset<flatbuffers::Vector<'a, flatbuffers::ForwardsUOffset<KeyVal<'a>>>>>,
    pub istio_version: Option<flatbuffers::WIPOffset<&'a str>>,
    pub mesh_id: Option<flatbuffers::WIPOffset<&'a str>>,
    pub app_containers: Option<flatbuffers::WIPOffset<flatbuffers::Vector<'a, flatbuffers::ForwardsUOffset<&'a str>>>>,
    pub cluster_id: Option<flatbuffers::WIPOffset<&'a str>>,
    pub instance_ips: Option<flatbuffers::WIPOffset<flatbuffers::Vector<'a, flatbuffers::ForwardsUOffset<&'a str>>>>,
}
impl<'a> Default for FlatNodeArgs<'a> {
  #[inline]
  fn default() -> Self {
    FlatNodeArgs {
      name: None,
      namespace: None,
      labels: None,
      owner: None,
      workload_name: None,
      platform_metadata: None,
      istio_version: None,
      mesh_id: None,
      app_containers: None,
      cluster_id: None,
      instance_ips: None,
    }
  }
}

pub struct FlatNodeBuilder<'a: 'b, 'b> {
  fbb_: &'b mut flatbuffers::FlatBufferBuilder<'a>,
  start_: flatbuffers::WIPOffset<flatbuffers::TableUnfinishedWIPOffset>,
}
impl<'a: 'b, 'b> FlatNodeBuilder<'a, 'b> {
  #[inline]
  pub fn add_name(&mut self, name: flatbuffers::WIPOffset<&'b  str>) {
    self.fbb_.push_slot_always::<flatbuffers::WIPOffset<_>>(FlatNode::VT_NAME, name);
  }
  #[inline]
  pub fn add_namespace(&mut self, namespace: flatbuffers::WIPOffset<&'b  str>) {
    self.fbb_.push_slot_always::<flatbuffers::WIPOffset<_>>(FlatNode::VT_NAMESPACE, namespace);
  }
  #[inline]
  pub fn add_labels(&mut self, labels: flatbuffers::WIPOffset<flatbuffers::Vector<'b , flatbuffers::ForwardsUOffset<KeyVal<'b >>>>) {
    self.fbb_.push_slot_always::<flatbuffers::WIPOffset<_>>(FlatNode::VT_LABELS, labels);
  }
  #[inline]
  pub fn add_owner(&mut self, owner: flatbuffers::WIPOffset<&'b  str>) {
    self.fbb_.push_slot_always::<flatbuffers::WIPOffset<_>>(FlatNode::VT_OWNER, owner);
  }
  #[inline]
  pub fn add_workload_name(&mut self, workload_name: flatbuffers::WIPOffset<&'b  str>) {
    self.fbb_.push_slot_always::<flatbuffers::WIPOffset<_>>(FlatNode::VT_WORKLOAD_NAME, workload_name);
  }
  #[inline]
  pub fn add_platform_metadata(&mut self, platform_metadata: flatbuffers::WIPOffset<flatbuffers::Vector<'b , flatbuffers::ForwardsUOffset<KeyVal<'b >>>>) {
    self.fbb_.push_slot_always::<flatbuffers::WIPOffset<_>>(FlatNode::VT_PLATFORM_METADATA, platform_metadata);
  }
  #[inline]
  pub fn add_istio_version(&mut self, istio_version: flatbuffers::WIPOffset<&'b  str>) {
    self.fbb_.push_slot_always::<flatbuffers::WIPOffset<_>>(FlatNode::VT_ISTIO_VERSION, istio_version);
  }
  #[inline]
  pub fn add_mesh_id(&mut self, mesh_id: flatbuffers::WIPOffset<&'b  str>) {
    self.fbb_.push_slot_always::<flatbuffers::WIPOffset<_>>(FlatNode::VT_MESH_ID, mesh_id);
  }
  #[inline]
  pub fn add_app_containers(&mut self, app_containers: flatbuffers::WIPOffset<flatbuffers::Vector<'b , flatbuffers::ForwardsUOffset<&'b  str>>>) {
    self.fbb_.push_slot_always::<flatbuffers::WIPOffset<_>>(FlatNode::VT_APP_CONTAINERS, app_containers);
  }
  #[inline]
  pub fn add_cluster_id(&mut self, cluster_id: flatbuffers::WIPOffset<&'b  str>) {
    self.fbb_.push_slot_always::<flatbuffers::WIPOffset<_>>(FlatNode::VT_CLUSTER_ID, cluster_id);
  }
  #[inline]
  pub fn add_instance_ips(&mut self, instance_ips: flatbuffers::WIPOffset<flatbuffers::Vector<'b , flatbuffers::ForwardsUOffset<&'b  str>>>) {
    self.fbb_.push_slot_always::<flatbuffers::WIPOffset<_>>(FlatNode::VT_INSTANCE_IPS, instance_ips);
  }
  #[inline]
  pub fn new(_fbb: &'b mut flatbuffers::FlatBufferBuilder<'a>) -> FlatNodeBuilder<'a, 'b> {
    let start = _fbb.start_table();
    FlatNodeBuilder {
      fbb_: _fbb,
      start_: start,
    }
  }
  #[inline]
  pub fn finish(self) -> flatbuffers::WIPOffset<FlatNode<'a>> {
    let o = self.fbb_.end_table(self.start_);
    flatbuffers::WIPOffset::new(o.value())
  }
}

impl core::fmt::Debug for FlatNode<'_> {
  fn fmt(&self, f: &mut core::fmt::Formatter<'_>) -> core::fmt::Result {
    let mut ds = f.debug_struct("FlatNode");
      ds.field("name", &self.name());
      ds.field("namespace", &self.namespace());
      ds.field("labels", &self.labels());
      ds.field("owner", &self.owner());
      ds.field("workload_name", &self.workload_name());
      ds.field("platform_metadata", &self.platform_metadata());
      ds.field("istio_version", &self.istio_version());
      ds.field("mesh_id", &self.mesh_id());
      ds.field("app_containers", &self.app_containers());
      ds.field("cluster_id", &self.cluster_id());
      ds.field("instance_ips", &self.instance_ips());
      ds.finish()
  }
}
#[inline]
#[deprecated(since="2.0.0", note="Deprecated in favor of `root_as...` methods.")]
pub fn get_root_as_flat_node<'a>(buf: &'a [u8]) -> FlatNode<'a> {
  unsafe { flatbuffers::root_unchecked::<FlatNode<'a>>(buf) }
}

#[inline]
#[deprecated(since="2.0.0", note="Deprecated in favor of `root_as...` methods.")]
pub fn get_size_prefixed_root_as_flat_node<'a>(buf: &'a [u8]) -> FlatNode<'a> {
  unsafe { flatbuffers::size_prefixed_root_unchecked::<FlatNode<'a>>(buf) }
}

#[inline]
/// Verifies that a buffer of bytes contains a `FlatNode`
/// and returns it.
/// Note that verification is still experimental and may not
/// catch every error, or be maximally performant. For the
/// previous, unchecked, behavior use
/// `root_as_flat_node_unchecked`.
pub fn root_as_flat_node(buf: &[u8]) -> Result<FlatNode, flatbuffers::InvalidFlatbuffer> {
  flatbuffers::root::<FlatNode>(buf)
}
#[inline]
/// Verifies that a buffer of bytes contains a size prefixed
/// `FlatNode` and returns it.
/// Note that verification is still experimental and may not
/// catch every error, or be maximally performant. For the
/// previous, unchecked, behavior use
/// `size_prefixed_root_as_flat_node_unchecked`.
pub fn size_prefixed_root_as_flat_node(buf: &[u8]) -> Result<FlatNode, flatbuffers::InvalidFlatbuffer> {
  flatbuffers::size_prefixed_root::<FlatNode>(buf)
}
#[inline]
/// Verifies, with the given options, that a buffer of bytes
/// contains a `FlatNode` and returns it.
/// Note that verification is still experimental and may not
/// catch every error, or be maximally performant. For the
/// previous, unchecked, behavior use
/// `root_as_flat_node_unchecked`.
pub fn root_as_flat_node_with_opts<'b, 'o>(
  opts: &'o flatbuffers::VerifierOptions,
  buf: &'b [u8],
) -> Result<FlatNode<'b>, flatbuffers::InvalidFlatbuffer> {
  flatbuffers::root_with_opts::<FlatNode<'b>>(opts, buf)
}
#[inline]
/// Verifies, with the given verifier options, that a buffer of
/// bytes contains a size prefixed `FlatNode` and returns
/// it. Note that verification is still experimental and may not
/// catch every error, or be maximally performant. For the
/// previous, unchecked, behavior use
/// `root_as_flat_node_unchecked`.
pub fn size_prefixed_root_as_flat_node_with_opts<'b, 'o>(
  opts: &'o flatbuffers::VerifierOptions,
  buf: &'b [u8],
) -> Result<FlatNode<'b>, flatbuffers::InvalidFlatbuffer> {
  flatbuffers::size_prefixed_root_with_opts::<FlatNode<'b>>(opts, buf)
}
#[inline]
/// Assumes, without verification, that a buffer of bytes contains a FlatNode and returns it.
/// # Safety
/// Callers must trust the given bytes do indeed contain a valid `FlatNode`.
pub unsafe fn root_as_flat_node_unchecked(buf: &[u8]) -> FlatNode {
  flatbuffers::root_unchecked::<FlatNode>(buf)
}
#[inline]
/// Assumes, without verification, that a buffer of bytes contains a size prefixed FlatNode and returns it.
/// # Safety
/// Callers must trust the given bytes do indeed contain a valid size prefixed `FlatNode`.
pub unsafe fn size_prefixed_root_as_flat_node_unchecked(buf: &[u8]) -> FlatNode {
  flatbuffers::size_prefixed_root_unchecked::<FlatNode>(buf)
}
#[inline]
pub fn finish_flat_node_buffer<'a, 'b>(
    fbb: &'b mut flatbuffers::FlatBufferBuilder<'a>,
    root: flatbuffers::WIPOffset<FlatNode<'a>>) {
  fbb.finish(root, None);
}

#[inline]
pub fn finish_size_prefixed_flat_node_buffer<'a, 'b>(fbb: &'b mut flatbuffers::FlatBufferBuilder<'a>, root: flatbuffers::WIPOffset<FlatNode<'a>>) {
  fbb.finish_size_prefixed(root, None);
}
}  // pub mod Common
}  // pub mod Wasm

