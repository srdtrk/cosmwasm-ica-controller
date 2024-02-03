"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[132],{5368:(e,t,s)=>{s.r(t),s.d(t,{assets:()=>i,contentTitle:()=>a,default:()=>h,frontMatter:()=>c,metadata:()=>o,toc:()=>l});var n=s(5893),r=s(1151);const c={title:"QueryMsg",sidebar_label:"QueryMsg",sidebar_position:3,slug:"/contract-api/query-msg"},a="QueryMsg",o={id:"contract-api/query-msg",title:"QueryMsg",description:"The QueryMsg is the message that allows you to see the state of the contract. Therefore, it is important to",source:"@site/docs/contract-api/03-query-msg.mdx",sourceDirName:"contract-api",slug:"/contract-api/query-msg",permalink:"/cw-ica-controller/main/contract-api/query-msg",draft:!1,unlisted:!1,editUrl:"https://github.com/srdtrk/cw-ica-controller/tree/feat/docusaurus-docs/docs/docs/contract-api/03-query-msg.mdx",tags:[],version:"current",sidebarPosition:3,frontMatter:{title:"QueryMsg",sidebar_label:"QueryMsg",sidebar_position:3,slug:"/contract-api/query-msg"},sidebar:"docsSidebar",previous:{title:"ExecuteMsg",permalink:"/cw-ica-controller/main/contract-api/execute-msg"},next:{title:"Callbacks",permalink:"/cw-ica-controller/main/contract-api/callbacks"}},i={},l=[{value:"<code>GetChannel</code>",id:"getchannel",level:2},{value:"<code>GetContractState</code>",id:"getcontractstate",level:2},{value:"<code>ica_info</code>",id:"ica_info",level:3},{value:"<code>allow_channel_open_init</code>",id:"allow_channel_open_init",level:3},{value:"<code>callback_address</code>",id:"callback_address",level:3},{value:"<code>Ownership</code>",id:"ownership",level:2}];function d(e){const t={a:"a",code:"code",h1:"h1",h2:"h2",h3:"h3",p:"p",pre:"pre",...(0,r.a)(),...e.components};return(0,n.jsxs)(n.Fragment,{children:[(0,n.jsx)(t.h1,{id:"querymsg",children:(0,n.jsx)(t.code,{children:"QueryMsg"})}),"\n",(0,n.jsxs)(t.p,{children:["The ",(0,n.jsx)(t.code,{children:"QueryMsg"})," is the message that allows you to see the state of the contract. Therefore, it is important to\nunderstand the state of the contract."]}),"\n",(0,n.jsx)(t.h2,{id:"getchannel",children:(0,n.jsx)(t.code,{children:"GetChannel"})}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-rust",metastring:"reference",children:"https://github.com/srdtrk/cw-ica-controller/blob/v0.4.0/src/types/msg.rs#L103-L105\n"})}),"\n",(0,n.jsx)(t.p,{children:"This message is used to query the state of the ICS-27 channel as tracked by the contract. It returns"}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-rust",metastring:"reference",children:"https://github.com/srdtrk/cw-ica-controller/blob/v0.4.0/src/types/state.rs#L156-L164\n"})}),"\n",(0,n.jsx)(t.h2,{id:"getcontractstate",children:(0,n.jsx)(t.code,{children:"GetContractState"})}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-rust",metastring:"reference",children:"https://github.com/srdtrk/cw-ica-controller/blob/v0.4.0/src/types/msg.rs#L106-L108\n"})}),"\n",(0,n.jsx)(t.p,{children:"This message is used to query the state of the contract. It returns"}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-rust",metastring:"reference",children:"https://github.com/srdtrk/cw-ica-controller/blob/v0.4.0/src/types/state.rs#L31-L45\n"})}),"\n",(0,n.jsxs)(t.p,{children:["Lets look at the fields of the ",(0,n.jsx)(t.code,{children:"ContractState"}),":"]}),"\n",(0,n.jsx)(t.h3,{id:"ica_info",children:(0,n.jsx)(t.code,{children:"ica_info"})}),"\n",(0,n.jsx)(t.p,{children:"This field will be empty if the channel handshake has not been completed. Otherwise, it will contain the\nfollowing information:"}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-rust",metastring:"reference",children:"https://github.com/srdtrk/cw-ica-controller/blob/v0.4.0/src/types/state.rs#L109-L115\n"})}),"\n",(0,n.jsx)(t.h3,{id:"allow_channel_open_init",children:(0,n.jsx)(t.code,{children:"allow_channel_open_init"})}),"\n",(0,n.jsx)(t.p,{children:"This field is used internally and should be ignored. It will be removed in a future version."}),"\n",(0,n.jsx)(t.h3,{id:"callback_address",children:(0,n.jsx)(t.code,{children:"callback_address"})}),"\n",(0,n.jsxs)(t.p,{children:["This is the contract address that the ",(0,n.jsx)(t.code,{children:"cw-ica-controller"})," contract will send callbacks to. If this field is empty,\nthen the contract will not send callbacks."]}),"\n",(0,n.jsx)(t.h2,{id:"ownership",children:(0,n.jsx)(t.code,{children:"Ownership"})}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-rust",metastring:"reference",children:"https://github.com/larry0x/cw-plus-plus/blob/ownable-v0.5.0/packages/ownable/derive/src/lib.rs#L142-L144\n"})}),"\n",(0,n.jsxs)(t.p,{children:["This message type is provided by the ",(0,n.jsx)(t.a,{href:"https://crates.io/crates/cw-ownable",children:"cw-ownable"})," crate. It allows to query\nthe ownership of the contract. It returns ",(0,n.jsx)(t.code,{children:"Ownership<String>"}),":"]}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-rust",metastring:"reference",children:"https://github.com/larry0x/cw-plus-plus/blob/ownable-v0.5.0/packages/ownable/src/lib.rs#L14-L29\n"})})]})}function h(e={}){const{wrapper:t}={...(0,r.a)(),...e.components};return t?(0,n.jsx)(t,{...e,children:(0,n.jsx)(d,{...e})}):d(e)}},1151:(e,t,s)=>{s.d(t,{Z:()=>o,a:()=>a});var n=s(7294);const r={},c=n.createContext(r);function a(e){const t=n.useContext(c);return n.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function o(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(r):e.components||r:a(e.components),n.createElement(c.Provider,{value:t},e.children)}}}]);