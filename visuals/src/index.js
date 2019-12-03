/* global fetch */
import React, { Component } from "react";
import { render } from "react-dom";
import { StaticMap } from "react-map-gl";
import DeckGL from "@deck.gl/react";
import { HeatmapLayer } from "@deck.gl/aggregation-layers";

// Set your mapbox token here
const MAPBOX_TOKEN = process.env.MapboxAccessToken; // eslint-disable-line

// Source data GeoJSON
const DATA_URL = "http://localhost:8082/visual"; // eslint-disable-line

const INITIAL_VIEW_STATE = {
  longitude: -100,
  latitude: 40.7,
  zoom: 3,
  maxZoom: 15,
  pitch: 0,
  bearing: 0
};

/* eslint-disable react/no-deprecated */
export class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      heatmaps: null,
      hoveredPath: null
    };
    this._onHover = this._onHover.bind(this);
    this._renderTooltip = this._renderTooltip.bind(this);

    this._recalculateHeatmaps(this.props.data);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.data !== this.props.data) {
      this._recalculateHeatmaps(nextProps.data);
    }
  }

  _onHover({ x, y, object }) {
    this.setState({ x, y, hovered: object });
  }

  _renderTooltip() {
    const { x, y, hovered } = this.state;
    return (
      hovered && (
        <div className="tooltip" style={{ left: x, top: y }}>
          {`${hovered.name}: ${hovered.value} products`}
        </div>
      )
    );
  }

  _recalculateHeatmaps(data) {
    if (!data) {
      return;
    }

    const heatmaps = data.map(factory => {
      return {
        coordinates: [factory.Lon, factory.Lat],
        name: factory.Name,
        weight: factory.Count
      };
    });

    this.setState({ heatmaps });
  }

  _renderLayers() {
    const { data } = this.props;

    return [
      new HeatmapLayer({
        id: "heatmapLayer",
        data: this.state.heatmaps,
        radiusPixels: 100,
        getPosition: d => d.coordinates,
        getWeight: d => d.weight,
        pickable: true,
        onHover: this._onHover
      })
    ];
  }

  render() {
    const { mapStyle = "mapbox://styles/mapbox/light-v9" } = this.props;

    return (
      <DeckGL
        layers={this._renderLayers()}
        initialViewState={INITIAL_VIEW_STATE}
        controller={true}
      >
        <StaticMap
          reuseMaps
          mapStyle={mapStyle}
          preventStyleDiffing={true}
          mapboxApiAccessToken={MAPBOX_TOKEN}
        />

        {this._renderTooltip}
      </DeckGL>
    );
  }
}

export function renderToDOM(container) {
  render(<App />, container);

  fetch(DATA_URL)
    .then(response => response.json())
    .then(resolved => {
      render(<App data={resolved} />, container);
    });
}
